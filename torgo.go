package torgo

import (
	"net/url"
	"os"
	"strings"

	"github.com/petegabriel/torgo/download"
)

func ParseMagnet(magnetUri string) (error, *Magnet) {
	m := new(Magnet)
	return m.Parse(magnetUri)
}

// ParseTorrent
func ParseTorrent(fp string) (*Torrent, error) {
	t := new(Torrent)
	//download file
	err := download.Download(fp)
	if err != nil {
		return nil, err
	}

	link, err := url.Parse(fp)
	if err != nil {
		return nil, err
	}

	var fn string
	if paths := strings.Split(link.Path, "/"); len(paths) > 0 {
		fn = paths[len(paths)-1]
	}

	f, err := os.Open("./../" + fn)
	if err != nil {
		return nil, err
	}
	return t.Parse(f)
}

type parselable interface {
	Parse(string) (error, interface{})
}
