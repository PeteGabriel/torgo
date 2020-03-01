package bittorrent

import (
	"bytes"
	"fmt"
	"github.com/petegabriel/torgo/peers"
	"log"
	"net"
	"time"
)

//PeerConnection is a tcp connection between a peer and us
type PeerConnection struct{
	con net.Conn
	peer peers.Peer
}

//Dial starts a tcp connection with the given peer.
func Dial(p peers.Peer) (*PeerConnection, error){
	con, err := net.DialTimeout("tcp", p.String(), 5*time.Second)
	if err != nil {
		return nil, err
	}
	return &PeerConnection{
		peer:p,
		con:con,
	}, nil
}

//Complete a two-way BitTorrent handshake
func (c *PeerConnection) DoHandshake(ih , pid []byte) (*Handshake, error) {

	hs := NewHandshake(ih, pid)
	_, err := c.con.Write(hs.Serialize())
	if err != nil {
		return nil, fmt.Errorf("Could not handshake with %s. Disconnecting\n", c.peer.IP)
	}

	hsr, err := Deserialize(c.con)
	if err != nil {
		return nil, fmt.Errorf("Could not handshake with %s. Disconnecting\n", c.peer.IP)
	}

	//verify good result
	if !bytes.Equal(hsr.InfoHash[:], ih[:]) {
		return nil, fmt.Errorf("expected infohash %x but got %x", hsr.InfoHash, ih)
	}

	log.Printf("Handshake successful with %s.\n", c.peer.IP)
	return hsr, nil
}


