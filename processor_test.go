package snitch

import (
	"sync"
	"testing"

	"github.com/quipo/statsd"
	"github.com/stretchr/testify/assert"
	"github.com/xjewer/snitch/lib/config"
)

func Test_ProcessorRun(t *testing.T) {
	var wg sync.WaitGroup
	a := assert.New(t)
	lines := make(chan *Line, 0)
	reader := NewNoopReader(lines)

	cfg := config.Source{
		Name: "test",
	}

	parser, err := NewParser(reader, statsd.NoopClient{}, cfg)
	a.Nil(err)
	p := NewProcessor(parser, reader)

	wg.Add(1)
	go func() {
		p.Run()
		wg.Done()
	}()

	lines <- NewLine("test", nil)
	a.Equal(ErrProcessorStopped, p.Close())
	//a.Equal(0, len(lines))
	wg.Wait()
	close(lines)
}
