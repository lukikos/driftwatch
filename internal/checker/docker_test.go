package checker

import (
	"testing"
)

func TestCheckDockerContainer_MissingField(t *testing.T) {
	drifted, msg, err := checkDockerContainer(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing container field")
	}
	if drifted {
		t.Error("expected drifted=false on error")
	}
	if msg != "" {
		t.Errorf("expected empty message, got %q", msg)
	}
}

func TestCheckDockerContainer_EmptyContainer(t *testing.T) {
	_, _, err := checkDockerContainer(map[string]string{"container": ""})
	if err == nil {
		t.Fatal("expected error for empty container name")
	}
}

func TestCheckDockerContainer_DefaultExpectedRunning(t *testing.T) {
	// When docker is not available the command fails and we treat it as drift.
	// We just verify the default expected_status logic doesn't error out.
	drifted, msg, err := checkDockerContainer(map[string]string{"container": "nonexistent_container_xyz"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Container won't be found, so drift is expected.
	if !drifted {
		t.Error("expected drift for nonexistent container")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckDockerContainer_CustomExpectedStatus(t *testing.T) {
	// Provide a custom expected_status; container won't exist so drift fires.
	drifted, msg, err := checkDockerContainer(map[string]string{
		"container":       "nonexistent_container_xyz",
		"expected_status": "exited",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Error("expected drift for nonexistent container")
	}
	_ = msg
}

func TestCheckDockerContainer_ViaChecker(t *testing.T) {
	c := New()
	drifted, msg, err := c.Check("docker_container", map[string]string{
		"container": "nonexistent_container_xyz",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Error("expected drift for nonexistent container")
	}
	_ = msg
}
