package torgo

import (
	"github.com/matryer/is"
	"testing"
)

func TestParseOfInvalidUri(t *testing.T){
	is := is.New(t)
	uri := "xt=urn:btih:0678589e05f322707fc82546af38040ebe8af963&dn=Trailblazers.UK.S01E11.Of.Madchester.HDTV.x264-LiNKLE%5Beztv%5D.mkv%5Beztv%5D&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A80&tr=udp%3A%2F%2Fglotorrents.pw%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"

	err, _ := Parse(uri)

	is.True(err != nil)
}

func TestParseMagnetUri(t *testing.T) {

	is := is.New(t)

	url := "magnet:?xt=urn:btih:0678589e05f322707fc82546af38040ebe8af963&dn=Trailblazers.UK.S01E11.Of.Madchester.HDTV.x264-LiNKLE%5Beztv%5D.mkv%5Beztv%5D&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A80&tr=udp%3A%2F%2Fglotorrents.pw%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"
	err, tor := Parse(url)

	is.NoErr(err)

	//assert tor fields.
	is.Equal(tor.Xt, "urn:btih:0678589e05f322707fc82546af38040ebe8af963")
	is.Equal(tor.DisplayName, "Trailblazers.UK.S01E11.Of.Madchester.HDTV.x264-LiNKLE%5Beztv%5D.mkv%5Beztv%5D")

}
