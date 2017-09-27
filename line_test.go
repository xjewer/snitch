package snitch_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xjewer/snitch"
)

const (
	TestLineOk    = "[%s]	200	192.168.1.1:80	200	0.036	https	POST	/test	/test	OK	hostname"
	TestLineError = "[%s]	200	192.168.1.1:80 : 127.0.0.1:8000	504 : 200	0.7 : 0.002	https	POST	/test	/test	Error	hostname"
)

func TestGetEntries(t *testing.T) {
	a := assert.New(t)
	tl := snitch.NewLine(TestLineOk, nil)
	tl.Split("\t")
	result := tl.GetEntries()
	a.Equal(11, len(result))
}

func TestGetEntry(t *testing.T) {
	a := assert.New(t)
	tl := snitch.NewLine(TestLineOk, nil)
	tl.Split("\t")
	result, err := tl.GetEntry(1)
	a.Equal("200", result)
	a.Nil(err)
}

func TestGetEntryError(t *testing.T) {
	a := assert.New(t)
	tl := snitch.NewLine(TestLineError, nil)
	tl.Split("\t")
	_, err := tl.GetEntry(20000)
	a.Equal(snitch.ErrOutboundIndex, err)
}

func TestGetEntryError2(t *testing.T) {
	a := assert.New(t)
	tl := snitch.NewLine(TestLineOk, nil)
	_, err := tl.GetEntry(1)
	a.Equal(snitch.ErrOutboundIndex, err)
}
