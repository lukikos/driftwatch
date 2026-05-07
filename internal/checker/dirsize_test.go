package checker

import (
	"os"
	"path/filepath"
	"testing"
)

func makeTestDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}
	}
	return dir
}

func TestCheckDirSize_MissingPath(t *testing.T) {
	_, _, err := checkDirSize(map[string]string{"max_bytes": "1000"})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestCheckDirSize_MissingMaxBytes(t *testing.T) {
	dir := makeTestDir(t, nil)
	_, _, err := checkDirSize(map[string]string{"path": dir})
	if err == nil {
		t.Fatal("expected error for missing max_bytes")
	}
}

func TestCheckDirSize_InvalidMaxBytes(t *testing.T) {
	dir := makeTestDir(t, nil)
	_, _, err := checkDirSize(map[string]string{"path": dir, "max_bytes": "notanumber"})
	if err == nil {
		t.Fatal("expected error for invalid max_bytes")
	}
}

func TestCheckDirSize_NoDrift(t *testing.T) {
	dir := makeTestDir(t, map[string]string{
		"a.txt": "hello",
		"b.txt": "world",
	})
	// "hello" + "world" = 10 bytes; allow 1000
	drifted, msg, err := checkDirSize(map[string]string{"path": dir, "max_bytes": "1000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Errorf("expected no drift, got: %s", msg)
	}
}

func TestCheckDirSize_Drift(t *testing.T) {
	dir := makeTestDir(t, map[string]string{
		"big.txt": "this content is definitely longer than five bytes",
	})
	// allow only 5 bytes — should drift
	drifted, msg, err := checkDirSize(map[string]string{"path": dir, "max_bytes": "5"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Errorf("expected drift, got: %s", msg)
	}
}

func TestCheckDirSize_ViaChecker(t *testing.T) {
	dir := makeTestDir(t, map[string]string{"file.txt": "abc"})
	c := New()
	_, err := c.Check("dirsize", map[string]string{"path": dir, "max_bytes": "9999"})
	if err != nil {
		t.Fatalf("unexpected error via checker: %v", err)
	}
}

func TestCheckDirSize_NonExistentPath(t *testing.T) {
	_, _, err := checkDirSize(map[string]string{"path": "/nonexistent/path/xyz", "max_bytes": "1000"})
	if err == nil {
		t.Fatal("expected error for non-existent path")
	}
}
