package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "driftwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, `
webhook_url: https://example.com/hook
interval: 30s
checks:
  - name: check-env
    type: env_var
    env_var: HOME
    expected: /root
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.WebhookURL != "https://example.com/hook" {
		t.Errorf("wrong webhook_url: %s", cfg.WebhookURL)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("wrong interval: %v", cfg.Interval)
	}
}

func TestLoad_DefaultInterval(t *testing.T) {
	path := writeTempConfig(t, `
webhook_url: https://example.com/hook
checks:
  - name: env-check
    type: env_var
    env_var: PATH
    expected: /usr/bin
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != defaultInterval {
		t.Errorf("expected default interval %v, got %v", defaultInterval, cfg.Interval)
	}
}

func TestLoad_MissingWebhookURL(t *testing.T) {
	path := writeTempConfig(t, `
checks:
  - name: env-check
    type: env_var
    env_var: PATH
    expected: /usr/bin
`)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for missing webhook_url")
	}
}

func TestLoad_MissingCheckName(t *testing.T) {
	path := writeTempConfig(t, `
webhook_url: https://example.com/hook
checks:
  - type: env_var
    env_var: PATH
    expected: /usr/bin
`)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for missing check name")
	}
}

func TestLoad_HTTPCheck_MissingURL(t *testing.T) {
	path := writeTempConfig(t, `
webhook_url: https://example.com/hook
checks:
  - name: http-check
    type: http_status
    expected_status: 200
`)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for http_status check missing url")
	}
}

func TestLoad_HTTPCheck_Valid(t *testing.T) {
	path := writeTempConfig(t, `
webhook_url: https://example.com/hook
checks:
  - name: http-check
    type: http_status
    url: https://example.com
    expected_status: 200
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Checks) != 1 {
		t.Errorf("expected 1 check, got %d", len(cfg.Checks))
	}
}
