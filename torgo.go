package torgo

import (
	"net/url"
	"os"
	"strings"

	"github.com/petegabriel/torgo/download"
)

//ParseMagnet parses a magnet uri
func ParseMagnet(uri string) (*Magnet, error) {
	m := new(Magnet)
	return m.Parse(uri)
}

// ParseTorrent parses a .torrent url
func ParseTorrent(fp string) (*Torrent, error) {
	t := new(Torrent)
	//download file
	if err := download.Download(fp); err != nil {
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
