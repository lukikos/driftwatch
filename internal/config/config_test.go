package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "driftwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	yaml := `
check_interval: 30s
webhook:
  url: https://hooks.example.com/alert
  timeout: 5s
checks:
  - name: nginx-config
    type: file
    target: /etc/nginx/nginx.conf
    params:
      expected_hash: abc123
`
	path := writeTempConfig(t, yaml)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CheckInterval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.CheckInterval)
	}
	if cfg.Webhook.URL != "https://hooks.example.com/alert" {
		t.Errorf("unexpected webhook URL: %s", cfg.Webhook.URL)
	}
	if len(cfg.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(cfg.Checks))
	}
	if cfg.Checks[0].Name != "nginx-config" {
		t.Errorf("unexpected check name: %s", cfg.Checks[0].Name)
	}
}

func TestLoad_DefaultInterval(t *testing.T) {
	yaml := `
webhook:
  url: https://hooks.example.com/alert
checks:
  - name: sshd
    type: process
    target: sshd
`
	path := writeTempConfig(t, yaml)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CheckInterval != 60*time.Second {
		t.Errorf("expected default 60s, got %v", cfg.CheckInterval)
	}
}

func TestLoad_MissingWebhookURL(t *testing.T) {
	yaml := `
checks:
  - name: sshd
    type: process
    target: sshd
`
	path := writeTempConfig(t, yaml)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestLoad_MissingCheckName(t *testing.T) {
	yaml := `
webhook:
  url: https://hooks.example.com/alert
checks:
  - type: process
    target: sshd
`
	path := writeTempConfig(t, yaml)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing check name")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
