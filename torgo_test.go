package torgo

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
)

const (
	torrent = "https://webtorrent.io/torrents/cosmos-laundromat.torrent"
	magnet  = "magnet:?xt=urn:btih:c9e15763f722f23e98a29decdfae341b98d53056&dn=Cosmos+Laundromat&tr=udp%3A%2F%2Fexplodie.org%3A6969&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Ftracker.empire-js.us%3A1337&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=wss%3A%2F%2Ftracker.btorrent.xyz&tr=wss%3A%2F%2Ftracker.fastcast.nz&tr=wss%3A%2F%2Ftracker.openwebtorrent.com&ws=https%3A%2F%2Fwebtorrent.io%2Ftorrents%2F&xs=https%3A%2F%2Fwebtorrent.io%2Ftorrents%2Fcosmos-laundromat.torrent"
)

func TestParseTorrentFile(t *testing.T) {
	is := is.New(t)
	tt, err := Parse("./resources/MPH.pdf.torrent")

	tor := tt.(*Torrent)

	is.NoErr(err)
	is.Equal(tor.Announce, "https://academictorrents.com/announce.php")
	is.Equal(tor.Name, "MoralPsychHandbook.pdf")
	var ih string = fmt.Sprintf("%x", tor.InfoHash[:])
	is.True(ih == "90493c18f577d24d5646c5075193bf57faabdcf6")
	is.Equal(tor.PieceLength, 16384)
	is.Equal(tor.Length, 150932)
}

func TestParseTorrentUri(t *testing.T) {
	is := is.New(t)
	tt, err := Parse(torrent)

	tor := tt.(*Torrent)

	is.NoErr(err)
	is.Equal(tor.Announce, "udp://tracker.leechers-paradise.org:6969")
	is.Equal(tor.Name, "Cosmos Laundromat")
	var ih string = fmt.Sprintf("%x", tor.InfoHash[:])
	is.True(ih == "0e6d3306f0d3826736854865771a26798b68b4eb")
	is.Equal(tor.PieceLength, 262144)
	is.Equal(tor.Length, 0)
}

func TestParseOfInvalidTorrentUri(t *testing.T) {
	is := is.New(t)
	err, _ := Parse(torrent)

	is.True(err != nil)
	os.Remove("./../cosmos-laundromat.torrent")
}

func TestParseOfInvalidMagnetUri(t *testing.T) {
	is := is.New(t)
	uri := strings.Replace(magnet, "magnet:?", "", 1)
	_, err := Parse(uri)

	is.True(err != nil)
}

func TestParseMagnetUri(t *testing.T) {

	is := is.New(t)

	m, err := Parse(magnet)

	mag := m.(*Magnet)

	is.NoErr(err)

	//assert tor fields.
	is.Equal(mag.Origin, magnet)
	is.Equal(mag.Xt, "urn:btih:c9e15763f722f23e98a29decdfae341b98d53056")
	is.Equal(mag.DisplayName, "Cosmos Laundromat")
	is.Equal(mag.Source, "https://webtorrent.io/torrents/cosmos-laundromat.torrent")

	is.Equal(len(mag.Trackers), 8)
	is.Equal(mag.Trackers[0], "udp://explodie.org:6969")
	is.Equal(mag.Trackers[1], "udp://tracker.coppersurfer.tk:6969")
	is.Equal(mag.Trackers[2], "udp://tracker.empire-js.us:1337")
	is.Equal(mag.Trackers[3], "udp://tracker.leechers-paradise.org:6969")
	is.Equal(mag.Trackers[4], "udp://tracker.opentrackr.org:1337")
	is.Equal(mag.Trackers[5], "wss://tracker.btorrent.xyz")
	is.Equal(mag.Trackers[6], "wss://tracker.fastcast.nz")
	is.Equal(mag.Trackers[7], "wss://tracker.openwebtorrent.com")
	is.Equal(mag.Seed, "https://webtorrent.io/torrents/")
	is.Equal(mag.Hash, "c9e15763f722f23e98a29decdfae341b98d53056")
}

func TestD(t *testing.T) {
	tt, _ := Parse("./resources/debian.torrent")

	tor := tt.(*Torrent)
	Download(tor, "")

}
