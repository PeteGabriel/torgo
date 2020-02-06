package download

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestDownloadOfTorFile(t *testing.T) {
	is := is.New(t)
	p := "https://webtorrent.io/torrents/cosmos-laundromat.torrent"
	err := Download(p)
	is.NoErr(err)

	_, err = os.Stat("./../cosmos-laundromat.torrent")
	is.NoErr(err)
	//clean env
	os.Remove("./../cosmos-laundromat.torrent")
}
