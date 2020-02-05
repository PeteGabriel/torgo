package torgo

import (
	"io"

	"github.com/jackpal/bencode-go"
)

func (*Torrent) Parse(r io.Reader) (*Torrent, error) {
	t := &Torrent{}
	err := bencode.Unmarshal(r, &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

type Torrent struct {
	Announce string `bencode:"announce"`
	Info     Info   `bencode:"info"`
}

type Info struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}
