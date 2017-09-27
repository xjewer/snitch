package snitch_test

import (
	"sync"
	"testing"

	"github.com/quipo/statsd"
	"github.com/stretchr/testify/assert"
	"github.com/xjewer/snitch"
	"github.com/xjewer/snitch/lib/config"
	"github.com/xjewer/snitch/lib/simplelog"
)

func Test_ProcessorRun(t *testing.T) {
	var wg sync.WaitGroup
	a := assert.New(t)
	lines := make(chan *snitch.Line, 0)
	reader := snitch.NewNoopReader(lines)

	cfg := config.Source{
		Name: "test",
	}

	l := &simplelog.NoopLogger{}
	parser, err := snitch.NewParser(reader, statsd.NoopClient{}, cfg)
	a.Nil(err)
	p := snitch.NewProcessor(parser, reader, l)

	wg.Add(1)
	go func() {
		p.Run()
		wg.Done()
	}()

	lines <- snitch.NewLine("test", nil)
	a.Equal(snitch.ErrProcessorStopped, p.Close())
	//a.Equal(0, len(lines))
	wg.Wait()
	close(lines)
}
