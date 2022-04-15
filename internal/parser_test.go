package internal

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Topface/snitch/internal/lib/config"
	"github.com/quipo/statsd"
	"github.com/stretchr/testify/assert"
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
			result, err := GetAmount(tc.str, tc.sep)
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
			result := GetLastMatch(tc.str, tc.sep)
			a.Equal(tc.expectation, result)
		})
	}
}

func Test_GetElementString(t *testing.T) {
	a := assert.New(t)
	l := NewLine("1,2,3", nil)
	l.Split(",")
	str, err := GetElementString(l, 2, ":", false)
	a.Nil(err)
	a.Equal("3", str)

	_, err = GetElementString(l, 10, ":", false)
	a.Equal(ErrOutboundIndex, err)
}

func Test_GetElementAmount(t *testing.T) {
	a := assert.New(t)
	l := NewLine("1,2,3", nil)
	l.Split(",")
	str, err := GetElementAmount(l, 2, ":")
	a.Nil(err)
	a.Equal(float32(3), str)

	_, err = GetElementAmount(l, 10, ":")
	a.Equal(ErrOutboundIndex, err)
}

func Test_HandleLine(t *testing.T) {
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
	runHandleLineTestcases(cfg, t)
}

func Test_HandleLine2(t *testing.T) {
	cfg := config.Source{
		Delimiter: "	",
		Keys: []config.Key{
			{
				Key:    "$3.$6",
				Count:  true,
				Timing: "$4",
			},
		},
	}

	runHandleLineTestcases(cfg, t)
}

func runHandleLineTestcases(cfg config.Source, t *testing.T) {
	type testCase struct {
		str string
		err error
	}

	a := assert.New(t)
	p, err := NewParser(NewNoopReader(nil), statsd.NoopClient{}, cfg)
	a.Nil(err)
	h, ok := p.(*Handler)
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
			err := h.HandleLine(NewLine(tc.str, nil))
			if tc.err != nil {
				a.EqualError(tc.err, err.Error(), "Should be equal")
			} else {
				a.Nil(err, "Should be nil error")
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
	h, ok := p.(*Handler)
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
			err := h.HandleLine(NewLine(tc.str, nil))
			if tc.err != nil {
				a.EqualError(tc.err, err.Error())
			} else {
				a.Nil(err)
			}
		})
	}
}

func TestSubstitute(t *testing.T) {
	type testCase struct {
		str         string
		expectation string
	}
	cases := []testCase{
		{
			"10.1.12.13",
			"10_1_12_13",
		},
		{
			"test",
			"test",
		},
		{
			"1.2.3.4.5.6.7.8.9",
			"1_2_3_4_5_6_7_8_9",
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			a := assert.New(t)
			result := SubstituteDots(tc.str)
			a.EqualValues(tc.expectation, result)
		})
	}
}
