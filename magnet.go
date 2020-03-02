package torgo

import (
	"errors"
	"log"
	"net/url"
	"strconv"
	"strings"
)

const (
	//Display Name is a filename to display to the user, for convenience
	DisplayName = "dn"
	//Exact Topic is a URN containing file hash
	ExactTopic  = "xt"
	//Exact Source is a P2P link identified by a content-hash
	ExactSource = "xs"
	//Tracker URL address for Bittorrent downloads
	Tracker     = "tr"
	//Web seed is the payload data served over https
	WebSeed     = "ws"
	//Exact length is the size in bytes
	ExactLength = "xl"
)

//Parse a magnet uri into a magnet struct
func ParseMagnet(uri string) (*Magnet, error) {

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
	m := Magnet{
		Origin: uri,
		Xt: parts[ExactTopic],
		DisplayName: parts[DisplayName],
		Trackers: trs,
		Source: parts[ExactSource],
		Hash: h,
		Seed: parts[WebSeed],
		Size: convert(parts[ExactLength]),
	}

	return &m, nil
}

func decode(src string) string {
	if s, err := url.QueryUnescape(src); err != nil {
		log.Println(err.Error())
		return ""
	} else {
		return s
	}
}

func convert(src string) int {
	if i, err := strconv.Atoi(src); err != nil {
		log.Println(err.Error())
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
