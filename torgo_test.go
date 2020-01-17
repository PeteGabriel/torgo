package torgo

import (
	"github.com/petegabriel/torgo"
	"testing"
)

func TestParseMagnetUri(t *testing.T) {

	url := ""
	err, tor := torgo.Parse(url)

	if err != nil {
		t.Errorf("Error parsing magnet link: %s", err)
		t.FailNow()
	}

	//assert tor fields.
	if tor.Pieces == 0 {
		t.FailNow()
	}

	if tor.PieceLen == 0 {
		t.FailNow()
	}

	if tor.Name == "" {
		t.FailNow()
	}

	if tor.Len == 0 {
		t.FailNow()
	}

	if tor.Creation == 0 {
		t.FailNow()
	}

	if tor.Hash == "" {
		t.FailNow()
	}
}
