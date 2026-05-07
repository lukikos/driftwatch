package checker

import (
	"os"
	"testing"
)

func TestCheckDiskUsage_MissingPath(t *testing.T) {
	_, _, err := checkDiskUsage(map[string]string{
		"max_percent": "90",
	})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestCheckDiskUsage_MissingMaxPercent(t *testing.T) {
	_, _, err := checkDiskUsage(map[string]string{
		"path": "/",
	})
	if err == nil {
		t.Fatal("expected error for missing max_percent")
	}
}

func TestCheckDiskUsage_InvalidMaxPercent(t *testing.T) {
	_, _, err := checkDiskUsage(map[string]string{
		"path":        "/",
		"max_percent": "notanumber",
	})
	if err == nil {
		t.Fatal("expected error for non-numeric max_percent")
	}
}

func TestCheckDiskUsage_OutOfRangeMaxPercent(t *testing.T) {
	for _, val := range []string{"0", "101", "-5"} {
		_, _, err := checkDiskUsage(map[string]string{
			"path":        "/",
			"max_percent": val,
		})
		if err == nil {
			t.Fatalf("expected error for max_percent=%s", val)
		}
	}
}

func TestCheckDiskUsage_NoDrift(t *testing.T) {
	// Using "/" with 100% threshold should never report drift.
	drift, msg, err := checkDiskUsage(map[string]string{
		"path":        "/",
		"max_percent": "100",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift with 100%% threshold, got msg: %s", msg)
	}
}

func TestCheckDiskUsage_Drift(t *testing.T) {
	// Using threshold of 0.001% should always trigger drift on any real filesystem.
	drift, msg, err := checkDiskUsage(map[string]string{
		"path":        "/",
		"max_percent": "0.001",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Errorf("expected drift with 0.001%% threshold, got msg: %s", msg)
	}
}

func TestCheckDiskUsage_InvalidPath(t *testing.T) {
	_, _, err := checkDiskUsage(map[string]string{
		"path":        "/nonexistent/path/xyz",
		"max_percent": "80",
	})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestCheckDiskUsage_ViaChecker(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "diskusage-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	c := New()
	drift, _, err := c.Check("disk_usage", map[string]string{
		"path":        tmpDir,
		"max_percent": "100",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift via checker with 100% threshold")
	}
}
