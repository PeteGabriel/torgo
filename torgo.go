package torgo

import (
	"errors"
	"strings"
)

func Parse(uri string) (error, *Torrent) {

	if !strings.Contains(uri, "magnet:?") {
		return errors.New("invalid magnet uri"), nil
	}

	return nil, nil
}

type Torrent struct {
	Origin 		string
	DisplayName string
	Hash        string
	Size        int64
	Xt 			string
	Addr 		string

	Creation    float64
	PieceLen    int
	Pieces      int
}
