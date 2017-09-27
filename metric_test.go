package snitch_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xjewer/snitch"
	"github.com/xjewer/snitch/lib/config"
)

func Test_makeMetrics(t *testing.T) {
	const Prefix = "test"

	type testCase struct {
		value       []config.Key
		err         error
		expectation []*snitch.Metric
	}
	cases := []testCase{
		{
			value: []config.Key{
				{Key: "key1.$1.$5", Count: true, Timing: "$2", Delimiter: "|  "},
			},
			err: nil,
			expectation: []*snitch.Metric{
				snitch.NewMetric([]snitch.KeyPath{
					snitch.NewKeyPath("test", 0, false),
					snitch.NewKeyPath("key1", 0, false),
					snitch.NewKeyPath("", 1, true),
					snitch.NewKeyPath("", 5, true),
				}, true,
					true,
					2,
					"|  ",
				),
			},
		},
		{
			value: []config.Key{
				{Key: "key2.path.test", Count: false, Timing: "$1", Delimiter: " "},
			},
			err: nil,
			expectation: []*snitch.Metric{
				snitch.NewMetric([]snitch.KeyPath{
					snitch.NewKeyPath("test", 0, false),
					snitch.NewKeyPath("key2", 0, false),
					snitch.NewKeyPath("path", 0, false),
					snitch.NewKeyPath("test", 0, false),
				}, false,
					true,
					1,
					" ",
				),
			},
		},
		{
			value: []config.Key{
				{Key: "key2.path.test", Count: false, Timing: "t", Delimiter: " "},
			},
			err:         snitch.ErrEmptyVarName,
			expectation: []*snitch.Metric{},
		},
		{
			value: []config.Key{
				{Key: "--$d", Count: true, Delimiter: " -- "},
			},
			err: nil,
			expectation: []*snitch.Metric{
				snitch.NewMetric([]snitch.KeyPath{
					snitch.NewKeyPath("test", 0, false),
					snitch.NewKeyPath("--$d", 0, false),
				}, true,
					false,
					0,
					" -- ",
				),
			},
		},
		{
			value: []config.Key{
				{Key: "key2.path.$y", Count: true, Delimiter: ", "},
			},
			err:         errors.New("strconv.Atoi: parsing \"y\": invalid syntax"),
			expectation: []*snitch.Metric{},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			a := assert.New(t)
			result, err := snitch.MakeMetrics(tc.value, Prefix)
			a.Equal(tc.expectation, result)
			if tc.err != nil {
				a.EqualError(tc.err, err.Error())
			} else {
				a.Nil(err)
			}
		})
	}
}

func Test_getVarName(t *testing.T) {
	type testCase struct {
		value       string
		err         error
		expectation int
	}
	cases := []testCase{
		{"$1", nil, 1},
		{"$", snitch.ErrEmptyVarName, 0},
		{"", snitch.ErrEmptyVarName, 0},
		{"$rr", errors.New("strconv.Atoi: parsing \"rr\": invalid syntax"), 0},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			a := assert.New(t)
			result, err := snitch.GetVarName(tc.value)
			a.Equal(tc.expectation, result)
			if tc.err != nil {
				a.EqualError(tc.err, err.Error())
			} else {
				a.Nil(err)
			}
		})
	}
}
