package snitch_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/quipo/statsd"
	"github.com/stretchr/testify/assert"
	"github.com/xjewer/snitch"
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
			snitch.ErrEmptyString,
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
			result, err := snitch.GetAmount(tc.str, tc.sep)
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
			result := snitch.GetLastMatch(tc.str, tc.sep)
			a.Equal(tc.expectation, result)
		})
	}
}

func Test_GetElementString(t *testing.T) {
	a := assert.New(t)
	l := snitch.NewLine("1,2,3", nil)
	l.Split(",")
	str, err := snitch.GetElementString(l, 2, ":", false)
	a.Nil(err)
	a.Equal("3", str)

	_, err = snitch.GetElementString(l, 10, ":", false)
	a.Equal(snitch.ErrOutboundIndex, err)
}

func Test_GetElementAmount(t *testing.T) {
	a := assert.New(t)
	l := snitch.NewLine("1,2,3", nil)
	l.Split(",")
	str, err := snitch.GetElementAmount(l, 2, ":")
	a.Nil(err)
	a.Equal(float32(3), str)

	_, err = snitch.GetElementAmount(l, 10, ":")
	a.Equal(snitch.ErrOutboundIndex, err)
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

	p, err := snitch.NewParser(snitch.NewNoopReader(nil), statsd.NoopClient{}, cfg)
	a.Nil(err)
	h, ok := p.(*snitch.Handler)
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
			snitch.ErrOutboundIndex,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			a := assert.New(t)
			err := h.HandleLine(snitch.NewLine(tc.str, nil))
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

	p, err := snitch.NewParser(snitch.NewNoopReader(nil), statsd.NoopClient{}, cfg)
	a.Nil(err)
	h, ok := p.(*snitch.Handler)
	a.True(ok)
	cases := []testCase{
		{
			"[22/Sep/2017:01:56:40 +0300]	200",
			snitch.ErrOutboundIndex,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			a := assert.New(t)
			err := h.HandleLine(snitch.NewLine(tc.str, nil))
			if tc.err != nil {
				a.EqualError(tc.err, err.Error())
			} else {
				a.Nil(err)
			}
		})
	}
}
