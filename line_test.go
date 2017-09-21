package snitch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TestLineOk    = "[%s]	200	192.168.1.1:80	200	0.036	https	POST	/test	/test	OK	hostname"
	TestLineError = "[%s]	200	192.168.1.1:80 : 127.0.0.1:8000	504 : 200	0.7 : 0.002	https	POST	/test	/test	Error	hostname"
)

func BenchmarkNewLine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tl := NewLine(TestLineError, nil)
		tl.Split("\t")
		tl.GetEntry(1)
		tl.GetEntry(5)
	}
}

func TestNewLine(t *testing.T) {
	a := assert.New(t)
	tl := NewLine(TestLineOk, nil)
	tl.Split("\t")
	result, err := tl.GetEntry(1)
	a.Equal(result, "200")
	a.Nil(err)
}

func TestNewLineError(t *testing.T) {
	a := assert.New(t)
	tl := NewLine(TestLineOk, nil)
	tl.Split("\t")
	_, err := tl.GetEntry(20000)
	a.Equal(err, ErrOutboundIndex)
}

func TestNewLineError2(t *testing.T) {
	a := assert.New(t)
	tl := NewLine(TestLineOk, nil)
	_, err := tl.GetEntry(1)
	a.Equal(err, ErrOutboundIndex)
}
