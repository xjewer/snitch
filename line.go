package snitch

import (
	"errors"
	"strings"
)

var (
	ErrEmptyString   = errors.New("empty string")
	ErrOutboundIndex = errors.New("outbound index")
)

// Line is a simplelog line structure
type Line struct {
	text    string
	err     error
	entries []string
}

// NewLine makes new text structure
func NewLine(t string, err error) *Line {
	return &Line{
		text: t,
		err:  err,
	}
}

func (l *Line) Split(sep string) {
	l.entries = strings.Split(l.text, sep)
}

func (l *Line) GetEntries() []string {
	return l.entries
}

func (l *Line) GetEntry(i int) (string, error) {
	if i >= len(l.entries) {
		return "", ErrOutboundIndex
	}

	return l.entries[i], nil
}
