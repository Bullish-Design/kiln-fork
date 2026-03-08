// @feature:cli Configuration file loading and type definitions
package config

import (
	"errors"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Theme             string `yaml:"theme"`
	Font              string `yaml:"font"`
	URL               string `yaml:"url"`
	Name              string `yaml:"name"`
	Input             string `yaml:"input"`
	Output            string `yaml:"output"`
	Mode              string `yaml:"mode"`
	Layout            string `yaml:"layout"`
	FlatURLs          bool   `yaml:"flat-urls"`
	DisableTOC        bool   `yaml:"disable-toc"`
	DisableLocalGraph bool   `yaml:"disable-local-graph"`
	Port              string `yaml:"port"`
	Log               string `yaml:"log"`
}

// Load reads a kiln.yaml file from the given path.
// Returns (nil, nil) if the file does not exist (not an error).
// Returns (*Config, nil) on success.
// Returns (nil, error) on parse failure.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		if errors.Is(err, io.EOF) {
			return &cfg, nil
		}
		return nil, err
	}
	return &cfg, nil
}
