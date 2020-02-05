package download

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

/*
Download torrent file from url.
The local file will assume the last element of the
segment path as its name.
*/
func Download(url string) error {

	fn, err := parseName(url)
	if err != nil {
		log.Println(err)
		return err
	}
	//Try to download
	r, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return errors.New("error downloading from " + url)
	}
	defer r.Body.Close()

	//Try to save info in locale file
	out, err := os.Create("./../" + fn)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, r.Body)
	return err
}

func parseName(u string) (string, error) {
	// Create the file
	// assume filename
	tmp := strings.Split(u, "/")
	var fn string
	if len(tmp) > 0 {
		fn = tmp[len(tmp)-1]
		return fn, nil
	} else {
		return "", errors.New("error found in url while parsing " + u)
	}

}
