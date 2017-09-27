package snitch_test

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xjewer/snitch"
	"github.com/xjewer/snitch/lib/config"
)

func Test_NewFileReaderWithZeroOffset(t *testing.T) {
	a := assert.New(t)
	const offsetFile = "./reader_test_offset1"

	// starts with second line
	err := cleanupOffsetFile(offsetFile, 0)
	a.Nil(err)

	expectations := []string{
		"[22/Sep/2017:01:56:25 +0300]	200	192.168.1.1:80 : 127.0.0.1:8000	504 : 200	0.7 : 0.002	https	POST	/test	/test	Error	hostname",
		"[22/Sep/2017:01:56:26 +0300]	200	192.168.1.1:80 : 127.0.0.1:8000	504 : 200	0.7 : 0.002	https	POST	/test	/test	Error	hostname",
		"[22/Sep/2017:01:56:39 +0300]	200	192.168.1.1:80	200	0.036	https	POST	/test	/test	OK	hostname",
		"[22/Sep/2017:01:56:27 +0300]	200	192.168.1.1:80 : 127.0.0.1:8000	504 : 200	0.7 : 0.002	https	POST	/test	/test	Error	hostname",
		"[22/Sep/2017:01:56:40 +0300]	200	192.168.1.1:80	200	0.036	https	POST	/test	/test	OK	hostname",
	}

	checkResults(a, expectations, offsetFile)
}

func Test_NewFileReaderWithOffset(t *testing.T) {
	a := assert.New(t)
	const offsetFile = "./reader_test_offset2"

	// starts with second line
	err := cleanupOffsetFile(offsetFile, 125)
	a.Nil(err)

	expectations := []string{
		"[22/Sep/2017:01:56:26 +0300]	200	192.168.1.1:80 : 127.0.0.1:8000	504 : 200	0.7 : 0.002	https	POST	/test	/test	Error	hostname",
		"[22/Sep/2017:01:56:39 +0300]	200	192.168.1.1:80	200	0.036	https	POST	/test	/test	OK	hostname",
		"[22/Sep/2017:01:56:27 +0300]	200	192.168.1.1:80 : 127.0.0.1:8000	504 : 200	0.7 : 0.002	https	POST	/test	/test	Error	hostname",
		"[22/Sep/2017:01:56:40 +0300]	200	192.168.1.1:80	200	0.036	https	POST	/test	/test	OK	hostname",
	}

	checkResults(a, expectations, offsetFile)
}

func checkResults(a *assert.Assertions, strs []string, offsetFile string) {
	var wg sync.WaitGroup

	cfg := config.Source{
		Name:       "test",
		ReOpen:     true,
		File:       "./reader_test.log",
		OffsetFile: offsetFile,
		MustExists: true,
	}

	reader, err := snitch.NewFileReader(cfg)
	a.Nil(err)

	lines := make(chan *snitch.Line)

	wg.Add(1)
	go func() {
		reader.GetLines(lines)
		wg.Done()
	}()

	read := func() {
		for _, e := range strs {
			select {
			case l, ok := <-lines:
				if !ok {
					fmt.Println("not ok")
					return
				}
				a.Equal(e, l.GetText(), "Strings should be equal")
			}
		}
	}

	read()

	reader.Close()
	close(lines)
	wg.Wait()

	os.Remove(offsetFile)
}

func cleanupOffsetFile(file string, offset int) error {
	// savePosition saves the last position to specific file
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(strconv.Itoa(offset))
	return nil
}
