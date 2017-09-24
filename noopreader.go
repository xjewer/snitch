package snitch

// noopReader allows use reader as a mock object without reading real source
// slice stings in lines returns in GetLines channel
type noopReader struct {
	//lines []*Line
	lines chan *Line
	done  chan struct{}
}

// NewNoopReader returns no-op LogReader struct
func NewNoopReader(lines chan *Line) LogReader {
	return &noopReader{
		lines: lines,
		done:  make(chan struct{}),
	}
}

func (r *noopReader) Close() error {
	close(r.done)
	return nil
}

func (r *noopReader) GetLines(lines chan<- *Line) {
	if r.lines == nil {
		return
	}
	for {
		select {
		case l := <-r.lines:
			lines <- l
		case <-r.done:
			return
		}
	}
}
