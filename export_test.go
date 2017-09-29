package snitch

import "github.com/xjewer/snitch/lib/config"

// line

func (l *Line) GetText() string {
	return l.text
}

// metric

type Metric = metric
type KeyPath = keyPath

func GetVarName(v string) (int, error) {
	return getVarName(v)
}

func MakeMetrics(keys []config.Key, prefix string) ([]*metric, error) {
	return makeMetrics(keys, prefix)
}

func NewKeyPath(v string, m int, isVar bool) KeyPath {
	return KeyPath{
		val:   v,
		match: m,
		isVar: isVar,
	}
}

func NewMetric(keys []KeyPath, count, timing bool, td int, d string) *Metric {
	return &metric{
		keyPaths:   keys,
		count:      count,
		timing:     timing,
		timingData: td,
		delimiter:  d,
	}
}

// parser
type Handler = handler

func GetElementAmount(l *Line, i int, sep string) (float32, error) {
	return getElementAmount(l, i, sep)
}

func GetAmount(s string, sep string) (float32, error) {
	return getAmount(s, sep)
}

func GetElementString(l *Line, i int, sep string, last bool) (string, error) {
	return getElementString(l, i, sep, last)
}

func GetLastMatch(s string, sep string) string {
	return getLastMatch(s, sep)
}

func SubstituteDots(s string) string {
	return substituteDots(s)
}
