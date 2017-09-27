package snitch

import (
	"gopkg.in/tomb.v1"
)

// noopReader allows use reader as a mock object without reading real source
// slice stings in lines returns in GetLines channel
type noopReader struct {
	tomb.Tomb
	lines chan *Line
}

// NewNoopReader returns no-op LogReader struct
func NewNoopReader(lines chan *Line) LogReader {
	return &noopReader{
		lines: lines,
	}
}

func (r *noopReader) GetName() string {
	return "noopReader"
}

func (r *noopReader) Close() error {
	r.Kill(ErrReaderIsFinished)
	return r.Wait()
}

func (r *noopReader) GetLines(lines chan<- *Line) {
	defer r.Done()
	if r.lines == nil {
		return
	}
	for {
		select {
		case l, ok := <-r.lines:
			if !ok {
				return
			}
			lines <- l
		case <-r.Dying():
			return
		}
	}
}
