package torgo

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"github.com/petegabriel/torgo/bittorrent"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/petegabriel/torgo/peers"
)



const (
	// MaxBlockSize is the largest number of bytes a request can ask for
	MaxBlockSize = 16384
	// MaxBacklog is the number of unfulfilled requests a client can have in its pipeline
	MaxBacklog = 5
)


/**
Download the given torrent.
*/
func Download(t interface{}, path string) ([]byte, error) {
	//ignore magnet uris for now
	tor := t.(*Torrent)

	peerID, err := genPeerID()
	if err != nil {
		return nil, err
	}

	prs, err := tor.RequestPeers(peerID[:])
	if err != nil {
		log.Println("Cannot request peers ")
		return nil, err
	}

	//good to go. Save info.
	tor.Peers = prs
	tor.PeerID = peerID[:]

	// Init queues for workers to retrieve work and send results
	workQueue := make(chan *piece, len(tor.PieceHashes))
	for index, hash := range tor.PieceHashes {
		length := tor.calculatePieceSize(index)
		workQueue <- &piece{index,hash,length}
	}

	results := make(chan *result)
	// Start workers
	for _, peer := range tor.Peers {
		go tor.startDownloadWorker(peer, workQueue, results)
	}

	// Collect results into a buffer until full
	buf := make([]byte, tor.Length)
	donePieces := 0
	for donePieces < len(tor.PieceHashes) {
		res := <-results
		begin, end := tor.calculateBoundsForPiece(res.index)
		copy(buf[begin:end], res.buf)
		donePieces++

		percent := float64(donePieces) / float64(len(tor.PieceHashes)) * 100
		numWorkers := runtime.NumGoroutine() - 1 // subtract 1 for main thread
		log.Printf("(%0.2f%%) Downloaded piece #%d from %d peers\n", percent, res.index, numWorkers)
	}

	close(workQueue)
	return buf, nil
}

func (tor *Torrent) startDownloadWorker(p peers.Peer, workQueue chan *piece, results chan *result) {
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
	if err = c.Unchoke(); err != nil {
		log.Println(err.Error())
		return
	}
	//send interested sign
	if err = c.Interested(); err != nil {
		log.Println(err.Error())
		return
	}

	for pieceOfWork := range workQueue {
		if !c.Bitfield.HasPiece(pieceOfWork.index) {
			workQueue <- pieceOfWork // Put piece back on the queue
			continue
		}

		// Download the piece
		buf, err := attemptDownloadPiece(c, pieceOfWork)
		if err != nil {
			log.Println("Exiting", err)
			workQueue <- pieceOfWork // Put piece back on the queue
			return
		}

		err = checkIntegrity(pieceOfWork, buf)
		if err != nil {
			log.Printf("piece #%d failed integrity check\n", pieceOfWork.index)
			workQueue <- pieceOfWork // Put piece back on the queue
			continue
		}

		c.SendHave(pieceOfWork.index)
		results <- &result{index: pieceOfWork.index, buf: buf}
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

func (t *Torrent) calculateBoundsForPiece(index int) (begin int, end int) {
	begin = index * t.PieceLength
	end = begin + t.PieceLength
	if end > t.Length {
		end = t.Length
	}
	return begin, end
}

func (t *Torrent) calculatePieceSize(index int) int {
	begin, end := t.calculateBoundsForPiece(index)
	return end - begin
}

func checkIntegrity(pieceOfWork *piece, buf []byte) error{
	hash := sha1.Sum(buf)
	if !bytes.Equal(hash[:], pieceOfWork.hash[:]) {
		return fmt.Errorf("Index %d failed integrity check", pieceOfWork.index)
	}
	return nil
}

func attemptDownloadPiece(c *bittorrent.PeerConnection, pow *piece) ([]byte, error){
	state := progress{
		index:  pow.index,
		client: c,
		buf:    make([]byte, pow.length),
	}

	// Setting a deadline helps get unresponsive peers unstuck.
	// 30 seconds is more than enough time to download a 262 KB piece
	c.Con.SetDeadline(time.Now().Add(30 * time.Second))
	defer c.Con.SetDeadline(time.Time{}) // Disable the deadline

	for state.downloaded < pow.length {
		// If unchoked, send requests until we have enough unfulfilled requests
		if !state.client.Choked {
			for state.backlog < MaxBacklog && state.requested < pow.length {
				blockSize := MaxBlockSize
				// Last block might be shorter than the typical block
				if pow.length-state.requested < blockSize {
					blockSize = pow.length - state.requested
				}

				err := c.SendRequest(pow.index, state.requested, blockSize)
				if err != nil {
					return nil, err
				}
				state.backlog++
				state.requested += blockSize
			}
		}

		err := state.readMessage()
		if err != nil {
			return nil, err
		}
	}

	return state.buf, nil
}

func (state *progress) readMessage() error {
	msg, err := state.client.Read() // this call blocks
	if err != nil {
		return err
	}

	if msg == nil { // keep-alive
		return nil
	}

	switch msg.ID {
	case bittorrent.UnchokeID:
		state.client.Choked = false
	case bittorrent.ChokeID:
		state.client.Choked = true
	case bittorrent.HaveID:
		index, err := bittorrent.ParseHave(msg)
		if err != nil {
			return err
		}
		state.client.Bitfield.SetPiece(index)
	case bittorrent.PieceID:
		n, err := bittorrent.ParsePiece(state.index, state.buf, msg)
		if err != nil {
			return err
		}
		state.downloaded += n
		state.backlog--
	}
	return nil
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

//piece represents some part to be download
type piece struct {
	index  int
	hash   [20]byte
	length int
}

type progress struct {
	index int
	client *bittorrent.PeerConnection
	buf        []byte
	downloaded int
	requested  int
	backlog    int
}

//piece represents some part that was downloaded
type result struct {
	index int
	buf   []byte
}
