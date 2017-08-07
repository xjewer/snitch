package ranger

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/hpcloud/tail"
)

type Reader struct {
	posFile  string
	done     chan struct{}
	noFollow bool
	tail     *tail.Tail
}

func New(f string, nf bool, positionFile string) (*Reader, error) {
	r := &Reader{
		done:     make(chan struct{}),
		noFollow: nf,
	}

	if positionFile != "" {
		r.posFile = positionFile
	}

	err := r.openTail(f)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Reader) Close() error {
	if r.posFile != "" {
		r.savePosition()
	}
	r.tail.Cleanup()
	r.tail.Kill(tail.ErrStop)
	return r.tail.Wait()
}

func (r *Reader) openTail(f string) error {
	config := tail.Config{
		ReOpen:    true,
		MustExist: true,
		Follow:    !r.noFollow,
	}

	if r.posFile != "" {
		offset, _ := r.getPosition()
		config.Location = &tail.SeekInfo{Offset: offset, Whence: io.SeekStart}
	} else {
		config.Location = &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}
	}

	t, err := tail.TailFile(f, config)
	if err != nil {
		return err
	}

	r.tail = t
	return nil
}

func (r *Reader) GetLines() <-chan *tail.Line {
	return r.tail.Lines
}

// savePosition save the last position to specific file
func (r *Reader) savePosition() error {
	f, err := os.OpenFile(r.posFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	offset, err := r.tail.Tell()
	if err != nil {
		return err
	}

	f.WriteString(strconv.Itoa(int(offset)))
	f.Close()
	return nil
}

// getPosition returns the last position of reading file
func (r *Reader) getPosition() (int64, error) {
	b, err := ioutil.ReadFile(r.posFile)
	if err != nil {
		return 0, err
	}

	if len(b) == 0 {
		return 0, nil
	}

	return strconv.ParseInt(string(b), 10, 0)
}
