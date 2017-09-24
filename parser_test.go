package snitch

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/quipo/statsd"
	"github.com/stretchr/testify/assert"
	"github.com/xjewer/snitch/lib/config"
)

func Test_GetAmount(t *testing.T) {
	type testCase struct {
		str         string
		sep         string
		err         error
		expectation float32
	}
	cases := []testCase{
		{
			"1, 2, 3, 4",
			", ",
			nil,
			10,
		},
		{
			"-",
			", ",
			ErrEmptyString,
			0,
		},
		{
			"1",
			", ",
			nil,
			1,
		},
		{
			"0.1 1.2 1.009 0.00 1",
			" ",
			nil,
			3.309,
		},
		{
			" e",
			"|",
			&strconv.NumError{Func: "ParseFloat", Num: " e", Err: strconv.ErrSyntax},
			0,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			a := assert.New(t)
			result, err := getAmount(tc.str, tc.sep)
			a.Equal(tc.expectation, result)
			if tc.err != nil {
				a.EqualError(tc.err, err.Error())
			} else {
				a.Nil(err)
			}
		})
	}
}

func Test_GetLastMatch(t *testing.T) {
	type testCase struct {
		str         string
		sep         string
		expectation string
	}
	cases := []testCase{
		{
			"test test2",
			" ",
			"test2",
		},
		{
			"test test2",
			"  ",
			"test test2",
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			a := assert.New(t)
			result := getLastMatch(tc.str, tc.sep)
			a.Equal(tc.expectation, result)
		})
	}
}

func Test_GetElementString(t *testing.T) {
	a := assert.New(t)
	l := NewLine("1,2,3", nil)
	l.Split(",")
	str, err := getElementString(l, 2, ":", false)
	a.Nil(err)
	a.Equal("3", str)

	_, err = getElementString(l, 10, ":", false)
	a.Equal(ErrOutboundIndex, err)
}

func Test_GetElementAmount(t *testing.T) {
	a := assert.New(t)
	l := NewLine("1,2,3", nil)
	l.Split(",")
	str, err := getElementAmount(l, 2, ":")
	a.Nil(err)
	a.Equal(float32(3), str)

	_, err = getElementAmount(l, 10, ":")
	a.Equal(ErrOutboundIndex, err)
}

func Test_Run(t *testing.T) {
	var wg sync.WaitGroup
	a := assert.New(t)
	lines := make(chan *Line, 0)
	reader := NewNoopReader(lines)

	cfg := config.Source{
		Name: "test",
	}

	p, err := NewParser(reader, statsd.NoopClient{}, cfg)
	a.Nil(err)

	wg.Add(1)
	go func() {
		p.Run()
		wg.Done()
	}()

	lines <- NewLine("test", nil)
	a.Nil(p.Close())
	a.Equal(0, len(lines))
	wg.Wait()
	close(lines)
}

func Test_HandleLine(t *testing.T) {
	type testCase struct {
		str string
		err error
	}

	a := assert.New(t)
	cfg := config.Source{
		Delimiter: "	",
		Keys: []config.Key{
			{
				Key:       "All.$3.$6",
				Count:     true,
				Timing:    "$4",
				Delimiter: " : ",
			},
		},
	}

	p, err := NewParser(NewNoopReader(nil), statsd.NoopClient{}, cfg)
	a.Nil(err)
	h, ok := p.(*handler)
	a.True(ok)
	cases := []testCase{
		{
			"[22/Sep/2017:01:56:40 +0300]	200	192.168.1.1:80	200	0.036	https	POST	/test	/test	OK	hostname",
			nil,
		},
		{
			"[22/Sep/2017:01:56:40 +0300]	200	192.168.1.1:80	200	-	https	POST	/test	/test	OK	hostname",
			nil,
		},
		{
			"[22/Sep/2017:01:56:40 +0300]	",
			ErrOutboundIndex,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			a := assert.New(t)
			err := h.handleLine(NewLine(tc.str, nil))
			if tc.err != nil {
				a.EqualError(tc.err, err.Error())
			} else {
				a.Nil(err)
			}
		})
	}
}

func Test_HandleLineError(t *testing.T) {
	type testCase struct {
		str string
		err error
	}

	a := assert.New(t)
	cfg := config.Source{
		Delimiter: "	",
		Keys: []config.Key{
			{
				Key:       "All.$1",
				Count:     true,
				Timing:    "$10",
				Delimiter: " - ",
			},
		},
	}

	p, err := NewParser(NewNoopReader(nil), statsd.NoopClient{}, cfg)
	a.Nil(err)
	h, ok := p.(*handler)
	a.True(ok)
	cases := []testCase{
		{
			"[22/Sep/2017:01:56:40 +0300]	200",
			ErrOutboundIndex,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			a := assert.New(t)
			err := h.handleLine(NewLine(tc.str, nil))
			if tc.err != nil {
				a.EqualError(tc.err, err.Error())
			} else {
				a.Nil(err)
			}
		})
	}
}
