package bittorrent

/**
Serialize the handshake structure into an array of bytes.

1. The length of the protocol identifier, which is always 19              -> 1 byte
2. The name of the protocol (BitTorrent protocol)
3. Eight reserved bytes, all set to 0                                     -> 8 bytes
4. The infohash that we calculated earlier to identify which file we want -> 20 bytes
5. The Peer ID that we made up to identify ourselves                      -> 20 bytes
                                                                             49 bytes
*/
func (h *Handshake) Serialize() []byte {

	knownBytes := 49
	buf := make([]byte, len(h.Pstr) + knownBytes)
	buf[0] = byte(len(h.Pstr))
	idx := 1
	idx += copy(buf[idx:], h.Pstr)
	idx += copy(buf[idx:], make([]byte, 8))
	idx += copy(buf[idx:], h.InfoHash[:])
	idx += copy(buf[idx:], h.PeerID[:])

	return buf
}

//Deserialize into an handshake structure.
func Deserialize() (*Handshake, error) {

	return nil, nil
}

// Handshake represent the tcp handshake between us and the tracker.
type Handshake struct {
	Pstr     string //protocol identifier which is always BitTorrent protocol
	InfoHash [20]byte
	PeerID   [20]byte //identify ourselves
}

func NewHandshake(infoHash, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr: "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID: peerID,
	}
}
