package snitch

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hpcloud/tail"
	"github.com/xjewer/snitch/lib/config"
)

var (
	ErrReaderIsFinished = errors.New("reader is finished")
)

type LogReader interface {
	Close() error
	GetLines(chan<- *Line)
	GetName() string
}

// fileReaders allows to read lines from file
// if offsetFile has been specified, reader keeps the last offset of file which was read till
// and will start since that
type fileReader struct {
	offsetFile string
	tail       *tail.Tail
	source     config.Source
}

// NewFileReader returns log reader from files
func NewFileReader(source config.Source) (LogReader, error) {
	r := &fileReader{
		source: source,
	}

	if source.OffsetFile != "" {
		r.offsetFile = source.OffsetFile
	}

	err := r.openFile(source.File)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// GetName returns reader name
func (r *fileReader) GetName() string {
	return r.source.Name
}

// CLose closes log reader
func (r *fileReader) Close() error {
	if r.offsetFile != "" {
		r.savePosition()
	}
	r.tail.Cleanup()
	r.tail.Kill(ErrReaderIsFinished)
	err := r.tail.Wait()
	if err != ErrReaderIsFinished {
		return err
	}

	return nil
}

// openFile opens log file to read its lines
func (r *fileReader) openFile(f string) error {
	c := tail.Config{
		ReOpen:    r.source.ReOpen,
		MustExist: r.source.MustExists,
		Follow:    !r.source.NoFollow,
	}

	if r.source.OffsetFile != "" {
		offset, err := r.getPosition()
		if err != nil {
			return err
		}
		c.Location = &tail.SeekInfo{Offset: offset, Whence: io.SeekStart}
	} else {
		c.Location = &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}
	}

	t, err := tail.TailFile(f, c)
	if err != nil {
		return err
	}

	r.tail = t
	return nil
}

// GetLines reads log lines and sends it to the channel
func (r *fileReader) GetLines(lines chan<- *Line) {
	for {
		select {
		case l, ok := <-r.tail.Lines:
			if !ok {
				log.Println("channel with lines has closed")
				return
			}

			if l == nil {
				// empty line
				log.Println("empty line")
				continue
			}
			lines <- NewLine(l.Text, l.Err)
		case <-r.tail.Dying():
			log.Println("dying")
			return
		}
	}
}

// savePosition saves the last position to specific file
func (r *fileReader) savePosition() error {
	f, err := os.OpenFile(r.offsetFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	offset, err := r.tail.Tell()
	if err != nil {
		return err
	}

	f.WriteString(strconv.Itoa(int(offset)))
	return nil
}

// getPosition returns the last position of reading file
func (r *fileReader) getPosition() (int64, error) {
	b, err := ioutil.ReadFile(r.source.OffsetFile)
	if err != nil {
		return 0, err
	}

	if len(b) == 0 {
		return 0, nil
	}

	return strconv.ParseInt(strings.Trim(string(b), "\n"), 10, 0)
}
