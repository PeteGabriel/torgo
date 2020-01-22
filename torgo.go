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

type Torrent struct {
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

func Parse(uri string) (error, *Torrent) {

	suf := "magnet:?"
	if !strings.Contains(uri, suf) {
		return errors.New("invalid magnet uri"), nil
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
			h = ih[len(ih) -1]
			parts[ExactTopic] = decode(tmp[1])
		} else {
			parts[tmp[0]] = decode(tmp[1])
		}
	}

	//parse TR links

	return nil, &Torrent{
		Origin:      uri,
		Xt:          parts[ExactTopic],
		DisplayName: parts[DisplayName],
		Trackers:    trs,
		Source:      parts[ExactSource],
		Hash:        h,
		Seed:        parts[WebSeed],
		Size:        convert(parts[ExactLength]),
	}
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

