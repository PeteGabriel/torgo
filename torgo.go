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

/**
Parse a magnet link or a .torrent file.
The .torrent file can be via url or a local file.
 */
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

//Parselable represents a type that can perform the parse operation
type Parselable interface {
	Parse(loc string) (Downloadable, error)
}

//Downloadable represents a type that can perform the utils operation
type Downloadable interface {
	Download() error
}
