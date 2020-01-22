package torgo

import (
	"github.com/matryer/is"
	"strings"
	"testing"
)

const magnet = "magnet:?xt=urn:btih:c9e15763f722f23e98a29decdfae341b98d53056&dn=Cosmos+Laundromat&tr=udp%3A%2F%2Fexplodie.org%3A6969&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Ftracker.empire-js.us%3A1337&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=wss%3A%2F%2Ftracker.btorrent.xyz&tr=wss%3A%2F%2Ftracker.fastcast.nz&tr=wss%3A%2F%2Ftracker.openwebtorrent.com&ws=https%3A%2F%2Fwebtorrent.io%2Ftorrents%2F&xs=https%3A%2F%2Fwebtorrent.io%2Ftorrents%2Fcosmos-laundromat.torrent"

func TestParseOfInvalidUri(t *testing.T){
	is := is.New(t)
	uri := strings.ReplaceAll(magnet, "magnet:?", "")
	err, _ := Parse(uri)

	is.True(err != nil)
}

func TestParseMagnetUri(t *testing.T) {

	is := is.New(t)

	err, tor := Parse(magnet)

	is.NoErr(err)

	//assert tor fields.
	is.Equal(tor.Origin, magnet)
	is.Equal(tor.Xt, "urn:btih:c9e15763f722f23e98a29decdfae341b98d53056")
	is.Equal(tor.DisplayName, "Cosmos Laundromat")
	is.Equal(tor.Hash, "c9e15763f722f23e98a29decdfae341b98d53056")
	// xs == https%3A%2F%2Fwebtorrent.io%2Ftorrents%2Fcosmos-laundromat.torrent

	is.True(len(tor.Addr) == 8)
	is.Equal(tor.Addr[0], "udp://explodie.org:6969")
	is.Equal(tor.Addr[1], "udp://tracker.coppersurfer.tk:6969")
	is.Equal(tor.Addr[2], "udp://tracker.empire-js.us:1337")
	is.Equal(tor.Addr[3], "udp://tracker.leechers-paradise.org:6969")
	is.Equal(tor.Addr[4], "udp://tracker.opentrackr.org:1337")
	is.Equal(tor.Addr[5], "wss://tracker.btorrent.xyz")
	is.Equal(tor.Addr[6], "wss://tracker.fastcast.nz")
	is.Equal(tor.Addr[7], "wss://tracker.openwebtorrent.com")

}
