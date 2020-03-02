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

	if len(ips) >= 1{
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

/*
func TestPeerConnection_DoHandshake(t *testing.T) {
	is := is2.New(t)
	var ih []byte
	var pid []byte

	c, err := Dial(*newPeer())
	is.NoErr(err)

	hs, err := c.DoHandshake(ih, pid)
	is.NoErr(err)

	is.Equal(hs.InfoHash, ih)
	is.Equal(hs.PeerID, pid)
	is.Equal(hs.Pstr, "BitTorrent protocol")
}
 */
