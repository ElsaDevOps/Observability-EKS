package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type ProviderConfig struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url"`
	APIKey string `yaml:"api_key"`
}

type Config struct {
	Providers []ProviderConfig `yaml:"providers"`
	Interval  time.Duration    `yaml:"interval"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
