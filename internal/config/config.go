package config

import (
	"fmt"
	"os"
	"strings"
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
	for i := range cfg.Providers {
		envName := fmt.Sprintf("%s_APIKEY", strings.ToUpper(cfg.Providers[i].Name))
		envValue := os.Getenv(envName)

		if envValue == "" {
			return nil, fmt.Errorf("environment variable %s not set", envName)
		}

		cfg.Providers[i].APIKey = envValue
	}

	return &cfg, nil
}
