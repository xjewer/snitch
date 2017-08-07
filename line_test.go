package snitch

import (
	"testing"
)

const (
	TestLineOk    = "[%s]	200	192.168.1.1:80	200	0.036	https	POST	/test	/test	OK	hostname"
	TestLineError = "[%s]	200	192.168.1.1:80 : 127.0.0.1:8000	504 : 200	0.7 : 0.002	https	POST	/test	/test	Error	hostname"
)

func BenchmarkNewLine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tl := TestLineError
		l := NewLine(tl, "\t")
		l.GetStatusHttpStatusCode()
		l.GetTiming()
		l.GetType()
	}
}
