package bittorrent

import (
	"bytes"
	"fmt"
	"github.com/petegabriel/torgo/peers"
	"github.com/petegabriel/torgo/peers/bitfield"
	"log"
	"net"
	"time"
)

//PeerConnection is a tcp connection between a Peer and us
type PeerConnection struct {
	Con  net.Conn
	Peer peers.Peer
	Bitfield bitfield.Bitfield
	Choked   bool
}

//Dial starts a tcp connection with the given Peer.
func Dial(p peers.Peer) (*PeerConnection, error) {
	con, err := net.DialTimeout("tcp", p.String(), 5*time.Second)
	if err != nil {
		return nil, err
	}
	return &PeerConnection{
		Peer: p,
		Con:  con,
	}, nil
}

func (p *PeerConnection) Read() (*Message, error) {
	content := make([]byte, 4)
	n, err := p.Con.Read(content)
	if err != nil{
		return nil, err
	}

	// keep-alive message
	if n == 0 {
		return nil, nil
	}

	m, err := DeserializeMsg(content)
	if err != nil || n == 0{
		return nil, err
	}
	return m, nil
}

//DoHandshake completes a two-way BitTorrent handshake
func (c *PeerConnection) DoHandshake(ih, pid []byte) (*Handshake, error) {

	hs := NewHandshake(ih, pid)
	_, err := c.Con.Write(hs.Serialize())
	if err != nil {
		return nil, fmt.Errorf("Could not handshake with %s. Disconnecting\n", c.Peer.IP)
	}

	hsr, err := Deserialize(c.Con)
	if err != nil {
		return nil, fmt.Errorf("Could not handshake with %s. Disconnecting\n", c.Peer.IP)
	}

	//verify good result
	if !bytes.Equal(hsr.InfoHash[:], ih[:]) {
		return nil, fmt.Errorf("expected infohash %x but got %x", hsr.InfoHash, ih)
	}

	log.Printf("Handshake successful with %s.\n", c.Peer.IP)
	return hsr, nil
}

/*
Unchoke Peer.
Choking is a temporary refusal to upload. It is
one of BitTorrent’s most powerful idea to deal
with those who only download but
never upload.
*/
func (c *PeerConnection) Unchoke() error {
	m := NewUnchokeMessage()
	if _, err := c.Con.Write(m.SerializeMsg()); err != nil {
		return err
	}
	return nil
}

//Interested in obtaining pieces of a file the Peer has
func (c *PeerConnection) Interested() error {
	m := NewInterestedMessage()
	if _, err := c.Con.Write(m.SerializeMsg()); err != nil {
		return err
	}
	return nil
}

//SendHave sends a Have message to the peer
func (c *PeerConnection) SendHave(idx int) error {
	m := NewHaveMessage(idx)
	if _, err := c.Con.Write(m.SerializeMsg()); err != nil {
		return err
	}
	return nil
}

// SendRequest sends a Request message to the peer
func (c *PeerConnection) SendRequest(index, begin, length int) error {
	req := FormatRequest(index, begin, length)
	_, err := c.Con.Write(req.SerializeMsg())
	return err
}