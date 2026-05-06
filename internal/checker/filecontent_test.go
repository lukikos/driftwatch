package checker

import (
	"os"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "driftwatch-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestCheckFileContent_MissingPath(t *testing.T) {
	_, _, err := checkFileContent(map[string]string{"contains": "foo"})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestCheckFileContent_MissingContainsAndMatches(t *testing.T) {
	path := writeTempFile(t, "hello world")
	_, _, err := checkFileContent(map[string]string{"path": path})
	if err == nil {
		t.Fatal("expected error when neither contains nor matches provided")
	}
}

func TestCheckFileContent_NoDrift_Contains(t *testing.T) {
	path := writeTempFile(t, "hello world")
	drift, msg, err := checkFileContent(map[string]string{"path": path, "contains": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift, got: %s", msg)
	}
}

func TestCheckFileContent_Drift_Contains(t *testing.T) {
	path := writeTempFile(t, "hello world")
	drift, msg, err := checkFileContent(map[string]string{"path": path, "contains": "goodbye"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Error("expected drift but got none")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckFileContent_NoDrift_Matches(t *testing.T) {
	path := writeTempFile(t, "version=1.2.3")
	drift, msg, err := checkFileContent(map[string]string{"path": path, "matches": `version=\d+\.\d+\.\d+`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift, got: %s", msg)
	}
}

func TestCheckFileContent_Drift_Matches(t *testing.T) {
	path := writeTempFile(t, "version=abc")
	drift, _, err := checkFileContent(map[string]string{"path": path, "matches": `version=\d+\.\d+\.\d+`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Error("expected drift but got none")
	}
}

func TestCheckFileContent_InvalidRegex(t *testing.T) {
	path := writeTempFile(t, "hello")
	_, _, err := checkFileContent(map[string]string{"path": path, "matches": `[invalid`})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestCheckFileContent_FileNotFound(t *testing.T) {
	_, _, err := checkFileContent(map[string]string{"path": "/nonexistent/file.txt", "contains": "foo"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCheckFileContent_ViaChecker(t *testing.T) {
	path := writeTempFile(t, "enabled=true")
	c := New()
	drift, msg, err := c.Check("file_content", map[string]string{"path": path, "contains": "enabled=false"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Error("expected drift")
	}
	if msg == "" {
		t.Error("expected non-empty message")
	}
}
