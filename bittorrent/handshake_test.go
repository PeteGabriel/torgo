package bittorrent

import (
	"crypto/rand"
	"strings"
	"testing"

	is2 "github.com/matryer/is"
)

func TestSerialize(t *testing.T) {
	is := is2.New(t)

	var ih = []byte("c9e15763f722f23e98a29decdfae341b98d53056")
	var peerId [20]byte
	rand.Read(peerId[:])
	hd := NewHandshake(ih[0:20], peerId[0:20])

	s := hd.Serialize()

	is.Equal(len(s), 19+49)
	is.Equal(string(s[28:48]), "c9e15763f722f23e98a2")
	is.Equal(string(s[48:]), string(peerId[:]))
}

func TestDeserialize(t *testing.T) {
	is := is2.New(t)
	r := strings.NewReader(string([]byte{ 19,66,105,116,84,111,114,114,101,110,116,32,
		112,114,111,116,111,99,111,108,0,0,0,0,0,0,0,0,99,57,
		101,49,53,55,54,51,102,55,50,50,102,50,51,101,57,56,
		97,50,14,199,216,199,253,84,126,176,180,105,136,165,
		51,212,42,228,88,67,154,19 }))

	dh, err := Deserialize(r)
	is.NoErr(err)

	is.Equal(dh.Pstr, "BitTorrent protocol")
	is.Equal(dh.InfoHash, []byte("c9e15763f722f23e98a2"))
	is.Equal(dh.PeerID, []byte{14,199,216,199,253,84,126,176,180,105,136,165,51,212,42,228,88,67,154,19})
}