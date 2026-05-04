package checker_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"testing"

	"github.com/yourusername/driftwatch/internal/checker"
	"github.com/yourusername/driftwatch/internal/config"
)

func TestCheckEnvVar_NoDrift(t *testing.T) {
	t.Setenv("MY_VAR", "expected-value")

	c := checker.New([]config.Check{
		{Name: "env-test", Type: "env_var", Target: "MY_VAR", Expected: "expected-value"},
	})
	results := c.RunAll()

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Drifted {
		t.Errorf("expected no drift, but drift was detected")
	}
	if results[0].Err != nil {
		t.Errorf("unexpected error: %v", results[0].Err)
	}
}

func TestCheckEnvVar_Drift(t *testing.T) {
	t.Setenv("MY_VAR", "actual-value")

	c := checker.New([]config.Check{
		{Name: "env-test", Type: "env_var", Target: "MY_VAR", Expected: "expected-value"},
	})
	results := c.RunAll()

	if !results[0].Drifted {
		t.Errorf("expected drift, but none detected")
	}
}

func TestCheckFileHash_NoDrift(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "driftwatch-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	content := []byte("stable content")
	f.Write(content)
	f.Close()

	expected := fmt.Sprintf("%x", sha256.Sum256(content))

	c := checker.New([]config.Check{
		{Name: "file-test", Type: "file_hash", Target: f.Name(), Expected: expected},
	})
	results := c.RunAll()

	if results[0].Drifted {
		t.Errorf("expected no drift, but drift was detected")
	}
}

func TestCheckFileHash_Drift(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "driftwatch-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("changed content"))
	f.Close()

	c := checker.New([]config.Check{
		{Name: "file-test", Type: "file_hash", Target: f.Name(), Expected: "deadbeef"},
	})
	results := c.RunAll()

	if !results[0].Drifted {
		t.Errorf("expected drift, but none detected")
	}
}

func TestUnknownCheckType(t *testing.T) {
	c := checker.New([]config.Check{
		{Name: "unknown", Type: "unsupported", Target: "x", Expected: "y"},
	})
	results := c.RunAll()

	if results[0].Err == nil {
		t.Errorf("expected error for unknown check type, got nil")
	}
}
