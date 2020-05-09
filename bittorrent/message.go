/*
All of the remaining messages in the protocol take the form of
<length prefix><message ID><payload>.

The length prefix is a four byte big-endian value.
The message ID is a single decimal byte.
The payload is message dependent.

Message 			<length prefix> <message ID> <payload>
Keep-alive 			0000 0 none
Choke 				0001 0 none
Unchoke 			0001 1 none
Interested 			0001 2 none
Not-interested 		0001 3 none
Have 				0005 4 Piece index
Bitfield 			0001+X 5 Bitfield
Request 			0013 6 <index><begin><length>
Piece 				0009+X 7 <index><begin><block>
Cancel 				0013 8 <index><begin><length>
port 				0003 9 <listen-port>
*/
package bittorrent

import (
	"encoding/binary"
	"fmt"
)

const (
	KeepAliveID = 0
	ChokeID = 0
	UnchokeID = 1
	InterestedID = 2
	NotInterestedID = 3
	HaveID = 4
	BitfieldID = 5
	RequestID = 6
	PieceID = 7
	CancelID = 8
	portID = 9
)

// NewInterestedMessage generates a create interested message
func NewInterestedMessage() *Message {
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, 1)
	return create(InterestedID, a[:4])
}

// NewHaveMessage generates a create have message
func NewHaveMessage(idx int) *Message {
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, uint32(idx))
	return create(HaveID, a[:4])
}

// NewUnchokeMessage generate a create unchoke message
func NewUnchokeMessage() *Message {
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, 1)
	return create(UnchokeID, a[:4])
}

//Message represents a message sent from client to Peer
type Message struct {
	Payload []byte
	ID      int
}

// FormatRequest creates a REQUEST message
func FormatRequest(index, begin, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return create(RequestID, payload)
}

// ParseHave parses a HAVE message
func ParseHave(msg *Message) (int, error) {
	if msg.ID != HaveID {
		return 0, fmt.Errorf("expected HAVE (ID %d), got ID %d", HaveID, msg.ID)
	}
	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("expected payload length 4, got length %d", len(msg.Payload))
	}
	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}

//ParsePiece parses a PIECE message
func ParsePiece(index int, buf []byte, msg *Message) (int, error){
	if msg.ID != PieceID {
		return 0, fmt.Errorf("expected HAVE (ID %d), got ID %d", PieceID, msg.ID)
	}
	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("payload length too short: %d", len(msg.Payload))
	}
	parsedIdx := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
	if parsedIdx != index {
		return 0, fmt.Errorf("expected index (%d), got index %d", index, parsedIdx)
	}
	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if begin >= len(buf) {
		return 0, fmt.Errorf("begin offset too high. %d >= %d", begin, len(buf))
	}
	data := msg.Payload[8:]
	if begin+len(data) > len(buf) {
		return 0, fmt.Errorf("data too long [%d] for offset %d with length %d", len(data), begin, len(buf))
	}

	copy(buf[begin:], data)
	return len(data), nil
}

// SerializeMsg a message
func (m *Message) SerializeMsg() []byte {
	length := uint32(len(m.Payload) + 1) // +1 for id
	buf := make([]byte, length)
	idx := 0
	idx += copy(buf[idx:], m.Payload)
	idx += copy(buf[idx:], []byte{byte(m.ID)})
	return buf
}

func DeserializeMsg(cnt []byte) (*Message, error){
	return &Message{
		ID:      int(cnt[0]),
		Payload: cnt[1:],
	}, nil
}


// NewMessage creates a create message representation
func create(id int, pref []byte) *Message {
	return &Message{
		Payload: pref[:4],
		ID:      id,
	}
}
