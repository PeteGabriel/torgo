package torgo

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

const (
	DisplayName = "dn"
	ExactTopic  = "xt"
	ExactSource = "xs"
	Tracker     = "tr"
	WebSeed     = "ws"
	ExactLength = "xl"
)

func (t *Magnet) Download() error {
	return nil
}

//Parse a magnet uri into a magnet struct
func (m *Magnet) Parse(uri string) (*Magnet, error) {

	suf := "magnet:?"
	if !strings.Contains(uri, suf) {
		return nil, errors.New("invalid magnet magnetUri")
	}

	//remove suffix
	ur := strings.Replace(uri, suf, "", len(suf))

	parts := make(map[string]string)
	trs := make([]string, 0)
	var h string

	urP := strings.Split(ur, "&")
	for _, p := range urP {
		tmp := strings.Split(p, "=")

		if tmp[0] == Tracker {
			trs = append(trs, decode(tmp[1]))
		} else if tmp[0] == ExactTopic {
			ih := strings.Split(tmp[1], ":")
			h = ih[len(ih)-1]
			parts[ExactTopic] = decode(tmp[1])
		} else {
			parts[tmp[0]] = decode(tmp[1])
		}
	}

	//parse TR links
	m.Origin = uri
	m.Xt = parts[ExactTopic]
	m.DisplayName = parts[DisplayName]
	m.Trackers = trs
	m.Source = parts[ExactSource]
	m.Hash = h
	m.Seed = parts[WebSeed]
	m.Size = convert(parts[ExactLength])

	return m, nil
}

func decode(src string) string {
	if s, err := url.QueryUnescape(src); err != nil {
		//TODO log error
		return ""
	} else {
		return s
	}
}

func convert(src string) int {
	if i, err := strconv.Atoi(src); err != nil {
		//TODO log error
		return -1
	} else {
		return i
	}
}

//A Magnet represent the info present in a magnet URI.
type Magnet struct {
	Origin      string   //original uri
	DisplayName string   //dn
	Hash        string   //info hash
	Size        int      // size in bytes
	Xt          string   // exact topic
	Trackers    []string //tr
	Source      string
	Seed        string

	Creation float64
	PieceLen int
	Pieces   int
}
