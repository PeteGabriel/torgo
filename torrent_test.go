package torgo

import (
	"testing"

	"github.com/matryer/is"
)

func TestTorrent_RequestPeers(t *testing.T) {
	is2 := is.New(t)

	pID, err := genPeerID()
	is2.NoErr(err)

	tor, err := ParseTor("./resources/debian.torrent")
	is2.NoErr(err)

	prs, err := tor.RequestPeers(pID[:20])
	is2.NoErr(err)

	is2.True(len(prs)>0)
}