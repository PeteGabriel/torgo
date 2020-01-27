package torgo

func ParseMagnet(magnetUri string) (error, *Magnet) {
	m := new(Magnet)
	return m.Parse(magnetUri)
}

func ParseTorrent(toUri string) (error, *Torrent) {
	t := new(Torrent)
	return t.Parse(toUri)
}


type parselable interface {
	Parse(string) (error, interface{})
}


