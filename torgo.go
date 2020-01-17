package torgo

import "fmt"

func Hello() {
	fmt.Println("Hello, World!")
}

type Torrent struct {
	Hash     string
	Len      int64
	Name     string
	Creation float64
	PieceLen int
	Pieces   int
}
