package internal

import (
	"sync"
	"testing"

	"github.com/Topface/snitch/internal/lib/config"
	"github.com/Topface/snitch/internal/lib/simplelog"
	"github.com/quipo/statsd"
	"github.com/stretchr/testify/assert"
)

func Test_ProcessorRun(t *testing.T) {
	var wg sync.WaitGroup
	a := assert.New(t)
	lines := make(chan *Line, 0)
	reader := NewNoopReader(lines)

	cfg := config.Source{
		Name: "test",
	}

	l := &simplelog.NoopLogger{}
	parser, err := NewParser(reader, statsd.NoopClient{}, cfg)
	a.Nil(err)
	p := NewProcessor(parser, reader, l)

	wg.Add(1)
	go func() {
		p.Run()
		wg.Done()
	}()

	lines <- NewLine("test", nil)
	a.Equal(ErrProcessorIsFinished, p.Close())
	a.Equal(0, len(lines))
	wg.Wait()
	close(lines)
}
