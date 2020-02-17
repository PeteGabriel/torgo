package torgo

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/petegabriel/torgo/download"
)


//Download the given torrent
func DownloadTorrent(tor *Torrent, path string) error {



	return nil
}

//Download the given torrent
func DownloadMagnet(mag *Magnet) error {
	return nil
}

//ParseMagnet parses a magnet uri
func ParseMagnet(uri string) (*Magnet, error) {
	m := new(Magnet)
	return m.Parse(uri)
}

// ParseTorrent parses a .torrent url
func ParseTorrent(fp string) (*Torrent, error) {
	t := new(Torrent)

	//Check if file exists locally
	f, err := os.Open(fp)
	if err != nil {
		log.Print("File not found locally.")
	}else {
		return t.Parse(f)
	}

	//download file
	if err := download.Download(fp); err != nil {
		log.Printf("error downloading file: %s", err.Error())
		return nil, err
	}

	link, err := url.Parse(fp)
	if err != nil {
		log.Printf("error parsing .torrent url: %s", err.Error())
		return nil, err
	}

	var fn string
	if paths := strings.Split(link.Path, "/"); len(paths) > 0 {
		fn = paths[len(paths)-1]
	}

	f, err = os.Open("./../" + fn)
	if err != nil {
		log.Printf("error opening file: %s", err.Error())
		return nil, err
	}
	return t.Parse(f)
}


