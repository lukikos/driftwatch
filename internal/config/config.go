package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level driftwatch configuration.
type Config struct {
	CheckInterval time.Duration `yaml:"check_interval"`
	Webhook       WebhookConfig `yaml:"webhook"`
	Checks        []CheckConfig `yaml:"checks"`
}

// WebhookConfig defines the alert destination.
type WebhookConfig struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
	Timeout time.Duration     `yaml:"timeout"`
}

// CheckConfig describes a single infrastructure check.
type CheckConfig struct {
	Name   string            `yaml:"name"`
	Type   string            `yaml:"type"`
	Target string            `yaml:"target"`
	Params map[string]string `yaml:"params"`
}

// Load reads and parses the YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.CheckInterval <= 0 {
		c.CheckInterval = 60 * time.Second
	}
	if c.Webhook.URL == "" {
		return fmt.Errorf("webhook.url is required")
	}
	if c.Webhook.Timeout <= 0 {
		c.Webhook.Timeout = 10 * time.Second
	}
	for i, ch := range c.Checks {
		if ch.Name == "" {
			return fmt.Errorf("checks[%d]: name is required", i)
		}
		if ch.Type == "" {
			return fmt.Errorf("checks[%d]: type is required", i)
		}
		if ch.Target == "" {
			return fmt.Errorf("checks[%d]: target is required", i)
		}
	}
	return nil
}
