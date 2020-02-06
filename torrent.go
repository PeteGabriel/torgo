package torgo

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/jackpal/bencode-go"
)

//Parse a .torrent file
func (*Torrent) Parse(r io.Reader) (*Torrent, error) {
	t := &tor{}
	err := bencode.Unmarshal(r, &t)
	if err != nil {
		return nil, err
	}

	return t.toTorrentFile()
}

//Torrent represents the info present in  a.torrent file
type Torrent struct {
	Announce    string
	InfoHash    [20]byte //hash from info struct
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func (t *tor) toTorrentFile() (*Torrent, error) {
	//hash info struct
	ih, err := t.Info.hash()
	if err != nil {
		return nil, err
	}

	//split pieces
	sp, err := t.Info.splitPieceHashes()
	if err != nil {
		return nil, err
	}

	tr := &Torrent{
		Announce:    t.Announce,
		Name:        t.Info.Name,
		Length:      t.Info.Length,
		PieceLength: t.Info.PieceLength,
		InfoHash:    ih,
		PieceHashes: sp,
	}
	return tr, nil
}

func (i *info) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20           // Length of SHA-1 hash
	buf := []byte(i.Pieces) //total amount of pieces (in bytes)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}

	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	//copy hashes in series of 20
	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil

}

func (i *info) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

type tor struct {
	Announce string `bencode:"announce"`
	Info     info   `bencode:"info"`
}

type info struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}
