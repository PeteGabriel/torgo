package torgo

import (
	"log"
	"strings"
)


/**
Download the given torrent.
*/
func Download(t Downloadable, path string) error {

	return nil
}

func Parse(loc string) (Downloadable, error) {
	suf := "magnet:?"
	if strings.Contains(loc, suf) {
		m := new(Magnet)
		d , err :=  m.Parse(loc)
		if err != nil {
			log.Print(err.Error())
		}
		return d, err
	}

	t := new(Torrent)
	d, err := t.Parse(loc)

	if err != nil {
		log.Print(err.Error())
	}
	return d, err
}

type Parselable interface {
	Parse(loc string) (Downloadable, error)
}

type Downloadable interface {
	Download() error
}
