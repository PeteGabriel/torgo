package torgo

import (
	"errors"
	"strings"
)

func Parse(uri string) (error, *Torrent) {

	suf := "magnet:?"
	if !strings.Contains(uri, suf) {
		return errors.New("invalid magnet uri"), nil
	}

	//remove suffix
	ur := strings.Replace(uri, suf, "", len(suf))

	parts := make(map[string]string)
	urP := strings.Split(ur, "&")
	for _, p := range urP {
		tmp := strings.Split(p, "=")
		parts[tmp[0]] = tmp[1]
	}

	return nil, &Torrent{
		Origin:      uri,
		Xt:          parts["xt"],
		DisplayName: parts["dn"],
		Addr:        parts["tr"],
	}
}

type Torrent struct {
	Origin      string //original uri
	DisplayName string //dn
	Hash        string //info hash
	Size        int64  // size in bytes
	Xt          string // exact topic
	Addr        string //tr

	Creation float64
	PieceLen int
	Pieces   int
}
