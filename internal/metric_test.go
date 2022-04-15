package internal

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Topface/snitch/internal/lib/config"
	"github.com/stretchr/testify/assert"
)

func Test_makeMetrics(t *testing.T) {
	const Prefix = "test"

	type testCase struct {
		value       []config.Key
		err         error
		expectation []*Metric
	}
	cases := []testCase{
		{
			value: []config.Key{
				{Key: "key1.$1.$5", Count: true, Timing: "$2", Delimiter: "|  "},
			},
			err: nil,
			expectation: []*Metric{
				NewMetric([]KeyPath{
					NewKeyPath("test", 0, false),
					NewKeyPath("key1", 0, false),
					NewKeyPath("", 1, true),
					NewKeyPath("", 5, true),
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
			expectation: []*Metric{
				NewMetric([]KeyPath{
					NewKeyPath("test", 0, false),
					NewKeyPath("key2", 0, false),
					NewKeyPath("path", 0, false),
					NewKeyPath("test", 0, false),
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
			err:         ErrEmptyVarName,
			expectation: []*Metric{},
		},
		{
			value: []config.Key{
				{Key: "--$d", Count: true, Delimiter: " -- "},
			},
			err: nil,
			expectation: []*Metric{
				NewMetric([]KeyPath{
					NewKeyPath("test", 0, false),
					NewKeyPath("--$d", 0, false),
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
			expectation: []*Metric{},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			a := assert.New(t)
			result, err := MakeMetrics(tc.value, Prefix)
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
		{"$", ErrEmptyVarName, 0},
		{"", ErrEmptyVarName, 0},
		{"$rr", errors.New("strconv.Atoi: parsing \"rr\": invalid syntax"), 0},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			a := assert.New(t)
			result, err := GetVarName(tc.value)
			a.Equal(tc.expectation, result)
			if tc.err != nil {
				a.EqualError(tc.err, err.Error())
			} else {
				a.Nil(err)
			}
		})
	}
}
