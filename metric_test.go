package snitch

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xjewer/snitch/lib/config"
)

func Test_makeMetrics(t *testing.T) {
	const Prefix = "test"

	type testCase struct {
		value       []config.Key
		err         error
		expectation []*metric
	}
	cases := []testCase{
		{
			value: []config.Key{
				{Key: "key1.$1.$5", Count: true, Timing: "$2", Delimiter: "|  "},
			},
			err: nil,
			expectation: []*metric{
				{
					path: []keyPath{
						{val: "test", match: 0, isVar: false},
						{val: "key1", match: 0, isVar: false},
						{match: 1, isVar: true},
						{match: 5, isVar: true},
					},
					count:      true,
					timing:     true,
					timingData: 2,
					delimiter:  "|  ",
				},
			},
		},
		{
			value: []config.Key{
				{Key: "key2.path.test", Count: false, Timing: "$1", Delimiter: " "},
			},
			err: nil,
			expectation: []*metric{
				{
					path: []keyPath{
						{val: "test", match: 0, isVar: false},
						{val: "key2", match: 0, isVar: false},
						{val: "path", match: 0, isVar: false},
						{val: "test", match: 0, isVar: false},
					},
					count:      false,
					timing:     true,
					timingData: 1,
					delimiter:  " ",
				},
			},
		},
		{
			value: []config.Key{
				{Key: "key2.path.test", Count: false, Timing: "t", Delimiter: " "},
			},
			err:         ErrEmptyVarName,
			expectation: []*metric{},
		},
		{
			value: []config.Key{
				{Key: "--$d", Count: true, Delimiter: " -- "},
			},
			err: nil,
			expectation: []*metric{
				{
					path: []keyPath{
						{val: "test", match: 0, isVar: false},
						{val: "--$d", match: 0, isVar: false},
					},
					count:     true,
					delimiter: " -- ",
				},
			},
		},
		{
			value: []config.Key{
				{Key: "key2.path.$y", Count: true, Delimiter: ", "},
			},
			err:         errors.New("strconv.Atoi: parsing \"y\": invalid syntax"),
			expectation: []*metric{},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			a := assert.New(t)
			result, err := makeMetrics(tc.value, Prefix)
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
			result, err := getVarName(tc.value)
			a.Equal(tc.expectation, result)
			if tc.err != nil {
				a.EqualError(tc.err, err.Error())
			} else {
				a.Nil(err)
			}
		})
	}
}
