package internal

import (
	"strconv"
	"strings"

	"github.com/Topface/snitch/internal/lib/config"
)

type keyPath struct {
	val   string
	match int
	isVar bool
}

type metric struct {
	keyPaths   []keyPath
	count      bool
	timing     bool
	timingData int
	delimiter  string
}

// makeMetrics makes metrics, that have to send to statsd with specific keys
func makeMetrics(keys []config.Key, prefix string) ([]*metric, error) {
	metrics := make([]*metric, 0)
	for _, k := range keys {
		m := &metric{keyPaths: make([]keyPath, 0), count: k.Count}
		if k.Timing != "" {
			td, err := getVarName(k.Timing)
			if err != nil {
				return metrics, err
			}
			m.timing = true
			m.timingData = td
		}

		m.keyPaths = append(m.keyPaths, keyPath{val: prefix})
		for _, p := range strings.Split(k.Key, ".") {
			if string(p[0]) == "$" {
				match, err := getVarName(p)
				if err != nil {
					return metrics, err
				}
				m.keyPaths = append(m.keyPaths, keyPath{isVar: true, match: match})
			} else {
				m.keyPaths = append(m.keyPaths, keyPath{val: p})
			}
		}
		m.delimiter = k.Delimiter

		metrics = append(metrics, m)
	}
	return metrics, nil
}

// getVarName returns a var name
func getVarName(v string) (int, error) {
	if len(v) <= 1 {
		return 0, ErrEmptyVarName
	}

	n, err := strconv.Atoi(v[1:])
	if err != nil {
		return 0, err
	}
	return n, nil
}
