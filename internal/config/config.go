// Package config loads and validates driftwatch configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const defaultInterval = 60 * time.Second

// Check represents a single drift check definition.
type Check struct {
	Name           string `yaml:"name"`
	Type           string `yaml:"type"`
	Path           string `yaml:"path,omitempty"`
	Expected       string `yaml:"expected,omitempty"`
	EnvVar         string `yaml:"env_var,omitempty"`
	URL            string `yaml:"url,omitempty"`
	ExpectedStatus int    `yaml:"expected_status,omitempty"`
}

// Config is the top-level configuration structure.
type Config struct {
	WebhookURL string        `yaml:"webhook_url"`
	Interval   time.Duration `yaml:"interval"`
	Checks     []Check       `yaml:"checks"`
	ServerAddr string        `yaml:"server_addr"`
	Cooldown   time.Duration `yaml:"cooldown"`
}

// Load reads and validates a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if cfg.Interval == 0 {
		cfg.Interval = defaultInterval
	}
	if cfg.ServerAddr == "" {
		cfg.ServerAddr = ":8080"
	}

	if cfg.WebhookURL == "" {
		return nil, errors.New("webhook_url is required")
	}
	for i, chk := range cfg.Checks {
		if chk.Name == "" {
			return nil, fmt.Errorf("check[%d]: name is required", i)
		}
		if chk.Type == "" {
			return nil, fmt.Errorf("check %q: type is required", chk.Name)
		}
		if chk.Type == "http_status" && chk.URL == "" {
			return nil, fmt.Errorf("check %q: url is required for http_status type", chk.Name)
		}
	}
	return &cfg, nil
}
