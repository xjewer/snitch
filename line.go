package ranger

import (
	"strconv"
	"strings"
)

const (
	TypeStatusOK = "200"
)

type line struct {
	entries []string
}

func NewLine(l string, separator string) line {
	return line{entries: strings.Split(l, separator)}
}

func (l *line) IsOk() bool {
	return l.entries[3] == TypeStatusOK
}

func (l *line) GetTiming() (int64, error) {

	f, err := strconv.ParseFloat(l.entries[4], 32)
	if err != nil {
		return 0, err
	}

	return int64(1000 * f), nil
}

func (l *line) GetType() string {
	return l.entries[6]
}
