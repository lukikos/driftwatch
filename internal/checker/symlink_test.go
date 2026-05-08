package checker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func TestCheckSymlink_MissingPath(t *testing.T) {
	_, _, err := checkSymlink(config.Check{Fields: map[string]string{}})
	if err == nil {
		t.Fatal("expected error for missing path field")
	}
}

func TestCheckSymlink_PathDoesNotExist(t *testing.T) {
	drift, msg, err := checkSymlink(config.Check{
		Fields: map[string]string{"path": "/nonexistent/symlink/path"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Errorf("expected drift when symlink does not exist, got msg: %s", msg)
	}
}

func TestCheckSymlink_PathIsNotSymlink(t *testing.T) {
	tmpDir := t.TempDir()
	realFile := filepath.Join(tmpDir, "realfile.txt")
	if err := os.WriteFile(realFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	drift, msg, err := checkSymlink(config.Check{
		Fields: map[string]string{"path": realFile},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Errorf("expected drift when path is not a symlink, got msg: %s", msg)
	}
}

func TestCheckSymlink_NoDrift_NoTargetCheck(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "target")
	link := filepath.Join(tmpDir, "link")
	if err := os.WriteFile(target, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	drift, _, err := checkSymlink(config.Check{
		Fields: map[string]string{"path": link},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift when symlink exists and no target check requested")
	}
}

func TestCheckSymlink_NoDrift_CorrectTarget(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "target")
	link := filepath.Join(tmpDir, "link")
	if err := os.WriteFile(target, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	drift, _, err := checkSymlink(config.Check{
		Fields: map[string]string{"path": link, "expected_target": target},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift when symlink resolves to correct target")
	}
}

func TestCheckSymlink_Drift_WrongTarget(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "target")
	wrongTarget := filepath.Join(tmpDir, "other")
	link := filepath.Join(tmpDir, "link")
	for _, f := range []string{target, wrongTarget} {
		if err := os.WriteFile(f, []byte("data"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(wrongTarget, link); err != nil {
		t.Fatal(err)
	}

	drift, msg, err := checkSymlink(config.Check{
		Fields: map[string]string{"path": link, "expected_target": target},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Errorf("expected drift for wrong symlink target, got msg: %s", msg)
	}
}

func TestCheckSymlink_ViaChecker(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "target")
	link := filepath.Join(tmpDir, "link")
	if err := os.WriteFile(target, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	c := New()
	drift, _, err := c.Check(config.Check{
		Name: "symlink-check",
		Type: "symlink",
		Fields: map[string]string{"path": link, "expected_target": target},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift via checker dispatch")
	}
}
