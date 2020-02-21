package torgo

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/jackpal/bencode-go"
	"github.com/petegabriel/torgo/download"
)



func (t *Torrent) Download() error {
	return nil
}

func (t *Torrent) Parse(loc string) (Downloadable, error){
	//Check if file exists locally
	f, err := os.Open(loc)
	if err != nil {
		log.Print("File not found locally.")
	}else {
		return t.parseReader(f)
	}

	//download file
	if err := download.Download(loc); err != nil {
		log.Printf("error downloading file: %s", err.Error())
		return nil, err
	}

	link, err := url.Parse(loc)
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
	return t.parseReader(f)
}


//Parse a .torrent file
func (*Torrent) parseReader(r io.Reader) (*Torrent, error) {
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
