package torgo

import (
	"crypto/rand"
	"github.com/petegabriel/torgo/bittorrent"
	"log"
	"strings"
	"sync"

	"github.com/petegabriel/torgo/peers"
)

/**
Download the given torrent.
*/
func Download(t interface{}, path string) (bool, error) {
	//ignore magnet uris for now
	tor := t.(*Torrent)

	peerID, err := genPeerID()
	if err != nil {
		return false, err
	}

	prs, err := tor.RequestPeers(peerID[:])
	if err != nil {
		log.Println("Cannot request peers ")
		return false, err
	}

	//good to go. Save info.
	tor.Peers = prs
	tor.PeerID = peerID[:]

	pieces := make(chan Piece, len(tor.Peers))
	results := make(chan Result)
	var wg sync.WaitGroup

	for idx, p := range tor.Peers {
		pieces <- Piece{
			index:  idx,
			length: tor.PieceLength,
			hash:   tor.PieceHashes[idx],
		}

		wg.Add(1)

		//start a worker
		go startDownloadWorker(&wg, p, tor, pieces, results)
	}


	//closer goroutine
	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		//TODO do something with the result pieces
		log.Printf("Result number %d has arrived\n", r.index)
	}

	return false, nil
}

func startDownloadWorker(wg *sync.WaitGroup, p peers.Peer, tor *Torrent, pieces chan Piece, res chan Result) {
	defer wg.Done()

	//connect to peer
	c, err := bittorrent.Dial(p)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//perform handshake
	_, err = c.DoHandshake(tor.InfoHash[:], tor.PeerID)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//send unchoke sign
	c.SendUnchoke()
	//send interested sign
	c.SendInterested()

	//get a piece
	piece := <-pieces

	//2 TODO process
	//3 notify progress
	res <- Result{
		index: piece.index,
	}
}

func genPeerID() ([20]byte, error) {
	var peerID [20]byte
	_, err := rand.Read(peerID[:])
	if err != nil {
		return [20]byte{}, err
	}
	return peerID, nil
}

/**
Parse a magnet link or a .torrent file.
The .torrent file can be via url or a local file.
*/
func Parse(loc string) (interface{}, error) {
	suf := "magnet:?"
	if strings.Contains(loc, suf) {
		m, err := ParseMagnet(loc)
		if err != nil {
			log.Print(err.Error())
		}
		return m, err
	}

	d, err := ParseTor(loc)

	if err != nil {
		log.Print(err.Error())
	}
	return d, err
}

//Parselable represents a type that can perform the parse operation
type Parselable interface {
	Parse(loc string) (interface{}, error)
}

//Piece represents some part to be download
type Piece struct {
	index  int
	hash   [20]byte
	length int
}

//Piece represents some part that was downloaded
type Result struct {
	index int
	buf   []byte
}
