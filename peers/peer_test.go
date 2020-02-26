package peers

import (
	"testing"

	is2 "github.com/matryer/is"
)

func TestUnmarshal(t *testing.T) {
	is := is2.New(t)

	ps := []byte{185,220,101,74,26,133,91,12,49,162,200,213,5,9,140,42,35,39,95,82,155,233,101,229,176,23,144,49,52,44,185,21,216,198,127,129,188,6,225,2,224,45,71,146,122,68,35,39,76,179,89,111,234,96,104,200,153,99,200,213,93,113,206,207,154,156,216,195,129,27,234,96,195,244,218,83,26,225,188,122,2,73,200,212}
	peers, err := Unmarshal(ps)

	is.NoErr(err)
	is.True(len(peers) == 13)
}

func TestUnmarshalWithEmptyInput(t *testing.T) {
	is := is2.New(t)

	var ps []byte
	peers, err := Unmarshal(ps)

	is.NoErr(err)
	is.True(len(peers) == 0)
}