package bittorrent

import (
	"log"
	"net"
	"testing"

	is2 "github.com/matryer/is"
	"github.com/petegabriel/torgo/peers"
)

func newPeer() *peers.Peer {
	p := peers.Peer{
		Port: 80,
	}
	const host = "www.google.com"
	ips, err := net.LookupIP(host)
	if err != nil {
		log.Fatal(err)
	}

	if len(ips) >= 1 {
		p.IP = ips[0]
	}
	return &p
}

func TestDial(t *testing.T) {
	is := is2.New(t)

	p := newPeer()
	c, err := Dial(*p)
	is.NoErr(err)
	is.True(c != nil)
}

func TestPeerConnection_Unchoke(t *testing.T) {
	is := is2.New(t)

	c, _ := Dial(*newPeer())
	err := c.Unchoke()
	is.NoErr(err)
}

func TestPeerConnection_Interested(t *testing.T) {
	is := is2.New(t)

	c, _ := Dial(*newPeer())
	err := c.Interested()
	is.NoErr(err)
}