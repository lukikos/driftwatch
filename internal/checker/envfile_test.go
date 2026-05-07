package checker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return p
}

func TestCheckEnvFile_NoDrift(t *testing.T) {
	p := writeTempEnvFile(t, "# comment\nAPP_ENV=production\nDEBUG=false\n")
	drift, msg, err := checkEnvFile(config.Check{
		Name:   "env-check",
		Fields: map[string]string{"path": p, "key": "APP_ENV", "expected": "production"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift, got message: %s", msg)
	}
}

func TestCheckEnvFile_Drift(t *testing.T) {
	p := writeTempEnvFile(t, "APP_ENV=staging\n")
	drift, msg, err := checkEnvFile(config.Check{
		Name:   "env-check",
		Fields: map[string]string{"path": p, "key": "APP_ENV", "expected": "production"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Fatal("expected drift but got none")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckEnvFile_KeyNotFound(t *testing.T) {
	p := writeTempEnvFile(t, "OTHER_KEY=value\n")
	drift, msg, err := checkEnvFile(config.Check{
		Name:   "env-check",
		Fields: map[string]string{"path": p, "key": "APP_ENV", "expected": "production"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Fatal("expected drift for missing key")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckEnvFile_MissingPath(t *testing.T) {
	_, _, err := checkEnvFile(config.Check{
		Name:   "env-check",
		Fields: map[string]string{"key": "APP_ENV", "expected": "production"},
	})
	if err == nil {
		t.Fatal("expected error for missing path field")
	}
}

func TestCheckEnvFile_MissingKey(t *testing.T) {
	p := writeTempEnvFile(t, "APP_ENV=production\n")
	_, _, err := checkEnvFile(config.Check{
		Name:   "env-check",
		Fields: map[string]string{"path": p, "expected": "production"},
	})
	if err == nil {
		t.Fatal("expected error for missing key field")
	}
}

func TestCheckEnvFile_FileNotFound(t *testing.T) {
	_, _, err := checkEnvFile(config.Check{
		Name:   "env-check",
		Fields: map[string]string{"path": "/nonexistent/.env", "key": "APP_ENV", "expected": "production"},
	})
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestCheckEnvFile_ViaChecker(t *testing.T) {
	p := writeTempEnvFile(t, "APP_ENV=production\n")
	ch := New()
	drift, _, err := ch.Run(config.Check{
		Name:   "env-check",
		Type:   "envfile",
		Fields: map[string]string{"path": p, "key": "APP_ENV", "expected": "production"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift")
	}
}
