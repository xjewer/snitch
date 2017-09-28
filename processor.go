package snitch

import (
	"github.com/xjewer/snitch/lib/simplelog"
	"gopkg.in/tomb.v1"
)

// Processor is a structure which gets the Reader, reads line by line from that and sends them to the Parser
type Processor struct {
	tomb.Tomb
	p Parser
	r LogReader
	l simplelog.Logger
}

// NewProcessor returns the New Processor
func NewProcessor(p Parser, r LogReader, l simplelog.Logger) *Processor {
	return &Processor{
		p: p,
		r: r,
		l: l,
	}
}

// Close processor and reader
func (p *Processor) Close() error {
	err := p.r.Close()
	if err != ErrReaderIsFinished {
		p.l.Println(err)
	}
	p.Kill(ErrProcessorIsFinished)
	return p.Wait()
}

// Run runs handler getting readers's simplelog lines and parse them
func (p *Processor) Run() {
	defer p.Done()
	lines := make(chan *Line, 0)
	defer close(lines)
	go p.r.GetLines(lines)
	for {
		select {
		case l := <-lines:
			if l.err != nil {
				p.l.Println("got line error", l.err)
				continue
			}

			err := p.p.HandleLine(l)
			if err != nil {
				p.l.Println(err)
				continue
			}
		case <-p.Dying():
			p.l.Printf("Closing %q ...\n", p.r.GetName())
			return
		}
	}
}
