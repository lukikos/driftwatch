package checker

import (
	"os"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func TestCheckFileExists_NoDrift_Exists(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "driftwatch-*")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	check := config.Check{
		Name: "tmp-file",
		Type: "file_exists",
		Fields: map[string]string{"path": f.Name(), "expected": "true"},
	}

	drift, msg, err := checkFileExists(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift, got message: %s", msg)
	}
}

func TestCheckFileExists_Drift_ShouldExistButMissing(t *testing.T) {
	check := config.Check{
		Name: "missing-file",
		Type: "file_exists",
		Fields: map[string]string{"path": "/tmp/driftwatch-nonexistent-xyz", "expected": "true"},
	}

	drift, msg, err := checkFileExists(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Error("expected drift, got none")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckFileExists_NoDrift_ShouldNotExist(t *testing.T) {
	check := config.Check{
		Name: "absent-file",
		Type: "file_exists",
		Fields: map[string]string{"path": "/tmp/driftwatch-nonexistent-xyz", "expected": "false"},
	}

	drift, _, err := checkFileExists(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift")
	}
}

func TestCheckFileExists_Drift_ShouldNotExistButPresent(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "driftwatch-*")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	check := config.Check{
		Name: "unexpected-file",
		Type: "file_exists",
		Fields: map[string]string{"path": f.Name(), "expected": "false"},
	}

	drift, msg, err := checkFileExists(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Error("expected drift, got none")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckFileExists_MissingPath(t *testing.T) {
	check := config.Check{
		Name:   "no-path",
		Type:   "file_exists",
		Fields: map[string]string{},
	}

	_, _, err := checkFileExists(check)
	if err == nil {
		t.Error("expected error for missing path field")
	}
}

func TestCheckFileExists_InvalidExpected(t *testing.T) {
	check := config.Check{
		Name:   "bad-expected",
		Type:   "file_exists",
		Fields: map[string]string{"path": "/tmp/anything", "expected": "yes"},
	}

	_, _, err := checkFileExists(check)
	if err == nil {
		t.Error("expected error for invalid 'expected' value")
	}
}

func TestCheckFileExists_DefaultExpectedTrue(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "driftwatch-*")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	// No 'expected' field — should default to "true" (file exists → no drift)
	check := config.Check{
		Name:   "default-expected",
		Type:   "file_exists",
		Fields: map[string]string{"path": f.Name()},
	}

	drift, _, err := checkFileExists(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift with default expected=true and existing file")
	}
}

func TestCheckFileExists_ViaChecker(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "driftwatch-*")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	check := config.Check{
		Name:   "via-checker",
		Type:   "file_exists",
		Fields: map[string]string{"path": f.Name()},
	}

	c := New()
	drift, _, err := c.Check(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift")
	}
}
