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

import "encoding/binary"

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

// Serialize a message
func (m *Message) Serialize() []byte {
	buf := make([]byte, 5)
	idx := 0
	idx += copy(buf[idx:], m.lenPrefix)
	idx += copy(buf[idx:], []byte{byte(m.mID)})
	return buf
}

func NewUnchokeMessage() *Message {
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, 1)
	return new(UnchokeID, a[:4])
}

//Message represents a message sent from client to peer
type Message struct {
	lenPrefix []byte
	mID       int
}

// NewMessage creates a new message representation
func new(id int, pref []byte) *Message {
	return &Message{
		lenPrefix: pref[:4],
		mID:       id,
	}
}
