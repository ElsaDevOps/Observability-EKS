package config

import (
	"os"
	"time"

	"github.com/ElsaDevOps/Observability-EKS/internal/provider"
	"gopkg.in/yaml.v2"
)

type Config struct {
	DefaultInterval time.Duration             `yaml:"default_interval"`
	Providers       []provider.ProviderConfig `yaml:"providers"`
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
