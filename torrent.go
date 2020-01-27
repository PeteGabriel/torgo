package torgo


func (*Torrent) Parse(uri string) (error, *Torrent){
	return nil, nil
}


type Torrent struct {
	Origin string
}