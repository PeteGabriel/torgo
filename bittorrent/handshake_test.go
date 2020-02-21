package bittorrent

import (
	"testing"

	is2 "github.com/matryer/is"
)

func TestSerialize(t *testing.T) {
	is := is2.New(t)

	var i = make([]byte, 20)
	var p = make([]byte, 20)
	h := NewHandshake(i[:], p[:])

}