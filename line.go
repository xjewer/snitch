package snitch

import (
	"strconv"
	"strings"
)

type line struct {
	entries []string
}

func NewLine(l string, separator string) line {
	return line{entries: strings.Split(l, separator)}
}

func (l *line) GetStatusHttpStatusCode() string {
	return getFirstMatch(l.entries[3])
}

func (l *line) GetTiming() (int64, error) {
	f, err := strconv.ParseFloat(getFirstMatch(l.entries[4]), 32)
	if err != nil {
		return 0, err
	}

	return int64(1000 * f), nil
}

func (l *line) GetType() string {
	return l.entries[6]
}

func getFirstMatch(s string) string {
	index := strings.Index(s, ":")
	if index >= 0 {
		return s[:index-1]
	}

	return s
}
