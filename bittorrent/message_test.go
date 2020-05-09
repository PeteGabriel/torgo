package bittorrent

import (
	is2 "github.com/matryer/is"
	"testing"
)

func TestNewUnchokeMessage(t *testing.T) {
	is := is2.New(t)
	m := NewUnchokeMessage()
	is.Equal(m.ID, 1)
	is.True(len(m.Payload) == 4)
}
func TestMessage_Serialize(t *testing.T) {
	is := is2.New(t)
	//0001 1
	m := NewUnchokeMessage()
	is.True(m != nil)

	s := m.SerializeMsg()
	is.True(s != nil)
	is.True(len(s) == 5)

	is.True(s[0] == 0)
	is.True(s[1] == 0)
	is.True(s[2] == 0)
	is.True(s[3] == 1)

	is.True(s[4] == 1)
}