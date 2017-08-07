package ranger

import (
	"github.com/quipo/statsd"
	"github.com/xjewer/ranger/lib/stats"
)

type template interface{}

type Parser interface {
	HandleLine(string) error
}

type handler struct {
	t template
	s statsd.Statsd
}

func NewParser(s statsd.Statsd) Parser {
	return &handler{s: s}
}

func (p *handler) HandleLine(s string) error {
	l := NewLine(s, "\t")

	t := l.GetType()

	stats.SendEvent(p.s, "All", t)

	if !l.IsOk() {
		stats.SendEvent(p.s, "Error", t)
	}

	timing, err := l.GetTiming()
	if err != nil {
		return err
	}

	stats.SendTiming(p.s, t, timing)

	return nil
}
