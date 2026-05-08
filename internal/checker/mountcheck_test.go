package checker

import (
	"os"
	"testing"
)

func TestCheckMountPoint_MissingPath(t *testing.T) {
	_, _, err := checkMountPoint(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestCheckMountPoint_InvalidExpected(t *testing.T) {
	_, _, err := checkMountPoint(map[string]string{
		"path":     "/tmp",
		"expected": "maybe",
	})
	if err == nil {
		t.Fatal("expected error for invalid expected value")
	}
}

func TestCheckMountPoint_NonExistentPath(t *testing.T) {
	_, _, err := checkMountPoint(map[string]string{
		"path": "/nonexistent/path/xyz",
	})
	if err == nil {
		t.Fatal("expected error for non-existent path")
	}
}

func TestCheckMountPoint_NotADirectory(t *testing.T) {
	f, err := os.CreateTemp("", "mountcheck-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Close()

	_, _, err = checkMountPoint(map[string]string{
		"path": f.Name(),
	})
	if err == nil {
		t.Fatal("expected error for non-directory path")
	}
}

func TestCheckMountPoint_DefaultExpectedMounted(t *testing.T) {
	// /tmp should exist on all platforms used in CI
	drift, detail, err := checkMountPoint(map[string]string{
		"path": "/tmp",
	})
	if err != nil {
		t.Skipf("skipping: %v", err)
	}
	// We don't assert drift value since it depends on the host;
	// we just ensure it runs without error and returns a detail string.
	_ = drift
	if detail == "" {
		t.Error("expected non-empty detail string")
	}
}

func TestCheckMountPoint_ViaChecker(t *testing.T) {
	c := New()
	drift, detail, err := c.Check("mount_point", map[string]string{
		"path": "/tmp",
	})
	if err != nil {
		t.Skipf("skipping: %v", err)
	}
	_ = drift
	if detail == "" {
		t.Error("expected non-empty detail string")
	}
}

func TestCheckMountPoint_UnknownTypeStillErrors(t *testing.T) {
	c := New()
	_, _, err := c.Check("mount_point_unknown", map[string]string{
		"path": "/tmp",
	})
	if err == nil {
		t.Fatal("expected error for unknown check type")
	}
}
