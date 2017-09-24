package snitch

import (
	"bytes"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/quipo/statsd"
	"github.com/xjewer/snitch/lib/config"
	"github.com/xjewer/snitch/lib/stats"
	"gopkg.in/tomb.v1"
)

var (
	ErrParserClose  = errors.New("parser close")
	ErrEmptyVarName = errors.New("empty var name")
)

// Parser parses log text from reader and sends statistics
type Parser interface {
	Run()
	Close() error
}

type handler struct {
	tomb.Tomb
	reader  LogReader
	statsd  statsd.Statsd
	metrics []*metric
	cfg     config.Source
}

// NewParser makes new Parser
func NewParser(r LogReader, s statsd.Statsd, cfg config.Source) (Parser, error) {
	m, err := makeMetrics(cfg.Keys, cfg.Prefix)
	if err != nil {
		return nil, err
	}
	return &handler{
		statsd:  s,
		reader:  r,
		metrics: m,
		cfg:     cfg,
	}, nil
}

// handleLine handles log text and sends statistics to statsd
func (h *handler) handleLine(l *Line) error {
	l.Split(h.cfg.Delimiter)
	for _, m := range h.metrics {
		key, err := makeKeyFromPaths(l, m)
		if err != nil {
			return err
		}
		if m.count {
			stats.SendEvent(h.statsd, key)
		}

		if m.timing {
			f, err := getElementAmount(l, m.timingData, m.delimiter)
			if err != nil && err != ErrEmptyString {
				return err
			}

			// statsd wants milliseconds
			stats.SendTiming(h.statsd, key, int64(1000*f))
		}
	}

	return nil
}

// Close reader and channels
func (h *handler) Close() error {
	h.Kill(ErrParserClose)
	h.Wait()
	return h.reader.Close()
}

// Run runs handler getting readers's log lines and parse them
func (h *handler) Run() {
	defer h.Done()
	lines := make(chan *Line, 0)
	defer close(lines)
	go h.reader.GetLines(lines)
	for {
		select {
		case l := <-lines:
			if l.err != nil {
				log.Println("got line error", l.err)
				continue
			}

			err := h.handleLine(l)
			if err != nil {
				log.Println(err)
				continue
			}
		case <-h.Dying():
			log.Printf("Closing %q ...\n", h.cfg.Name)
			return
		}
	}
}

// makeKeyFromPaths makes statsd key from keyPath
// key path sets the order and type of each sequences
func makeKeyFromPaths(l *Line, m *metric) (string, error) {
	//todo use sync.Pool
	var buffer bytes.Buffer
	for i, k := range m.path {
		if i != 0 {
			buffer.WriteString(".")
		}
		if k.isVar {
			m, err := getElementString(l, k.match, m.delimiter, true)
			if err != nil {
				return "", err
			}
			buffer.WriteString(m)
		} else {
			buffer.WriteString(k.val)
		}
	}

	return buffer.String(), nil
}

// getElementString returns specific entry from entries, id last value pass,
// it returns the last from chain of ", "
func getElementString(l *Line, i int, sep string, last bool) (string, error) {
	c, err := l.GetEntry(i)
	if err != nil {
		return "", err
	}

	if last {
		return getLastMatch(c, sep), nil
	}

	return c, nil
}

// getElementAmount returns amount of values from specific entry
func getElementAmount(l *Line, i int, sep string) (float32, error) {
	c, err := l.GetEntry(i)
	if err != nil {
		return 0.0, err
	}

	return getAmount(c, sep)
}

// getLastMatch returns the last columns after string separator
func getLastMatch(s string, sep string) string {
	index := strings.LastIndex(s, sep)
	if index >= 0 && index+1 < len(s) {
		return s[index+len(sep):]
	}

	return s
}

// getAmount returns amount of numbers in column, if error happen - returns it
func getAmount(s string, sep string) (float32, error) {
	var result float32
	if s == "-" {
		return result, ErrEmptyString
	}

	columns := strings.Split(s, sep)
	for _, c := range columns {
		f, err := strconv.ParseFloat(c, 32)
		if err != nil {
			return result, err
		}

		if f == 0 {
			// avoid needless addition
			continue
		}

		result += float32(f)
	}
	return result, nil
}
