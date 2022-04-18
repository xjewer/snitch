package stats

import (
	"time"

	"github.com/quipo/statsd"
)

// NewStatsd return the new statsd client
func NewStatsd(s string, prefix string, buffer int) statsd.Statsd {
	if s == "" {
		return statsd.NewStdoutClient("", prefix)
	}

	sc := statsd.NewStatsdClient(s, prefix)
	if buffer != 0 {
		sb := statsd.NewStatsdBuffer(time.Second*time.Duration(buffer), sc)
		return sb
	}
	return sc
}

// SendTiming writes timings in milliseconds
func SendTiming(s statsd.Statsd, key string, t int64) {
	s.Timing(key, t)
}

// SendEvent writes event types
func SendEvent(s statsd.Statsd, key string) {
	s.Incr(key, 1)
}
