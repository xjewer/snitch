package internal

import (
	"testing"

	"github.com/Topface/snitch/internal/lib/config"
	"github.com/quipo/statsd"
)

func Benchmark_parser(b *testing.B) {
	lines := make(chan *Line, 0)
	reader := NewNoopReader(lines)
	cfg := config.Source{
		Name:      "test",
		Delimiter: "\t",
		Keys: []config.Key{
			{Key: "prefix.$3.$6", Count: true, Timing: "$4", Delimiter: ", "},
		},
	}

	l := NewLine("[22/Sep/2017:01:56:40 +0300]	200	192.168.1.1:80	200	0.036	https	POST	/test	/test	OK	hostname", nil)

	p, err := NewParser(reader, statsd.NoopClient{}, cfg)
	if err != nil {
		b.Fatal(err)
	}

	h, ok := p.(*Handler)
	if !ok {
		b.Fatalf("wrong interface assertion")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := h.HandleLine(l)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}
