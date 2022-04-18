package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	// DefaultConfigPath is a default config path
	DefaultConfigPath = "config.yml"
)

var (
	// ErrConfigFileMissing is a missing config file error
	ErrConfigFileMissing = errors.New("config file is missing")
)

// Data is a main config structure
type Data struct {
	Sources []Source `yaml:"sources"`
}

// Source describes log parser
type Source struct {
	Name       string `yaml:"source"`
	File       string `yaml:"file"`
	Template   string `yaml:"template"`
	Delimiter  string `yaml:"delimiter"`
	Prefix     string `yaml:"prefix"`
	Keys       []Key  `yaml:"keys"`
	NoFollow   bool   `yaml:"noFollow"`
	MustExists bool   `yaml:"mustExists"`
	ReOpen     bool   `yaml:"reOpen"`
	OffsetFile string `yaml:"offsetFile"`
}

// Key describes statsd keys and their metrics
type Key struct {
	Key       string `yaml:"key"`
	Count     bool   `yaml:"count"`
	Timing    string `yaml:"timing"`
	Delimiter string `yaml:"delimiter"`
}

// Parse parses config file
func Parse(path string) (*Data, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, ErrConfigFileMissing
	}

	var data Data
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		return nil, fmt.Errorf("data parse error: %s", err)
	}

	return &data, nil
}
