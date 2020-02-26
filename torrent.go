package torgo

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackpal/bencode-go"
	"github.com/petegabriel/torgo/bittorrent"
	"github.com/petegabriel/torgo/peers"
	"github.com/petegabriel/torgo/utils"
)


const Port = 6881


func (t *Torrent) Download() error {
	var peerId [20]byte
	_, err := rand.Read(peerId[:])
	if err != nil {
		return err
	}

	peers, err := t.requestPeers(peerId[:])
	if err != nil {
		fmt.Println("Cannot request peers ")
		return err
	}

	t.Peers = peers
	t.PeerID = peerId[:]

	t.downloadFromPeers()

	return nil
}

/**
1. Start a TCP connection with the peer
2. Complete a two-way BitTorrent handshake
3. Ask to utils pieces
 */
func (t *Torrent) downloadFromPeers() {

	for _, p := range t.Peers{
		con, err := net.DialTimeout("tcp", p.String(), 5*time.Second)
		if err != nil {
			fmt.Println(err.Error())
		}
		if con != nil {
			defer con.Close()
		}

		hs := bittorrent.NewHandshake(t.InfoHash[:], t.PeerID)
		_, err = con.Write(hs.Serialize())
		if err != nil {
			log.Printf("Could not handshake with %s. Disconnecting\n", p.IP)
			return
		}

		hsr, err := bittorrent.Deserialize(con)
		if err != nil {
			log.Printf("Could not handshake with %s. Disconnecting\n", p.IP)
			return
		}

		fmt.Println(hsr)
	}

}

func (t *Torrent) Parse(loc string) (Downloadable, error){
	//Check if file exists locally
	f, err := os.Open(loc)
	if err != nil {
		log.Print("File not found locally.")
	}else {
		return t.parseReader(f)
	}

	//utils file
	if err := utils.DownTorFile(loc); err != nil {
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

func (t *Torrent) requestPeers(peerId []byte) ([]peers.Peer, error) {

	tu, err := t.getUrlTracker(peerId[:], Port)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	c := &http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(tu)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tracker := tracker{}
	err = bencode.Unmarshal(resp.Body, &tracker)
	if err != nil {
		return nil, err
	}

	return peers.Unmarshal([]byte(tracker.Peers))
}

//peerId identifies us when meeting the tracker.
func (t *Torrent) getUrlTracker(peerId []byte, port int) (string, error) {
	u, err := url.Parse(t.Announce)
	if err != nil {

	}

	//extra query params
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerId[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	u.RawQuery = params.Encode()

	return u.String(), nil
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
	Peers       []peers.Peer
	PeerID      []byte
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

type tracker struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}
