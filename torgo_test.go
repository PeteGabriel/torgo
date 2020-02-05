package torgo

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

const (
	torrent = "https://webtorrent.io/torrents/cosmos-laundromat.torrent"
	magnet  = "magnet:?xt=urn:btih:c9e15763f722f23e98a29decdfae341b98d53056&dn=Cosmos+Laundromat&tr=udp%3A%2F%2Fexplodie.org%3A6969&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Ftracker.empire-js.us%3A1337&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=wss%3A%2F%2Ftracker.btorrent.xyz&tr=wss%3A%2F%2Ftracker.fastcast.nz&tr=wss%3A%2F%2Ftracker.openwebtorrent.com&ws=https%3A%2F%2Fwebtorrent.io%2Ftorrents%2F&xs=https%3A%2F%2Fwebtorrent.io%2Ftorrents%2Fcosmos-laundromat.torrent"
)

func TestParseTorrentUri(t *testing.T) {
	is := is.New(t)
	tor, err := ParseTorrent(torrent)

	is.NoErr(err)
	is.Equal(tor.Announce, "udp://tracker.leechers-paradise.org:6969")
}

func TestParseOfInvalidTorrentUri(t *testing.T) {
	is := is.New(t)
	err, _ := ParseTorrent(torrent)

	is.True(err != nil)
}

func TestParseOfInvalidMagnetUri(t *testing.T) {
	is := is.New(t)
	uri := strings.Replace(magnet, "magnet:?", "", 1)
	err, _ := ParseMagnet(uri)

	is.True(err != nil)
}

func TestParseMagnetUri(t *testing.T) {

	is := is.New(t)

	err, tor := ParseMagnet(magnet)

	is.NoErr(err)

	//assert tor fields.
	is.Equal(tor.Origin, magnet)
	is.Equal(tor.Xt, "urn:btih:c9e15763f722f23e98a29decdfae341b98d53056")
	is.Equal(tor.DisplayName, "Cosmos Laundromat")
	is.Equal(tor.Source, "https://webtorrent.io/torrents/cosmos-laundromat.torrent")

	is.Equal(len(tor.Trackers), 8)
	is.Equal(tor.Trackers[0], "udp://explodie.org:6969")
	is.Equal(tor.Trackers[1], "udp://tracker.coppersurfer.tk:6969")
	is.Equal(tor.Trackers[2], "udp://tracker.empire-js.us:1337")
	is.Equal(tor.Trackers[3], "udp://tracker.leechers-paradise.org:6969")
	is.Equal(tor.Trackers[4], "udp://tracker.opentrackr.org:1337")
	is.Equal(tor.Trackers[5], "wss://tracker.btorrent.xyz")
	is.Equal(tor.Trackers[6], "wss://tracker.fastcast.nz")
	is.Equal(tor.Trackers[7], "wss://tracker.openwebtorrent.com")
	is.Equal(tor.Seed, "https://webtorrent.io/torrents/")
	is.Equal(tor.Hash, "c9e15763f722f23e98a29decdfae341b98d53056")
}
