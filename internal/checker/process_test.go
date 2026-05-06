package checker

import (
	"runtime"
	"testing"
)

func TestCheckProcessRunning_MissingField(t *testing.T) {
	_, _, err := checkProcessRunning(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing 'process' field, got nil")
	}
}

func TestCheckProcessRunning_EmptyProcess(t *testing.T) {
	_, _, err := checkProcessRunning(map[string]string{"process": ""})
	if err == nil {
		t.Fatal("expected error for empty 'process' field, got nil")
	}
}

func TestCheckProcessRunning_NoDrift(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping pgrep-based test on Windows")
	}

	// Use a process we know is running in any Unix-like CI environment.
	// "sh" or the test binary itself via its parent are typical choices;
	// we rely on the fact that the Go test runner spawns at least one shell.
	// A safer bet: look for the process that is definitely present.
	var processName string
	switch runtime.GOOS {
	case "darwin":
		processName = "launchd"
	default:
		processName = "init" // PID 1 on most Linux systems; may be "systemd"
	}

	drifted, msg, err := checkProcessRunning(map[string]string{"process": processName})
	if err != nil {
		// Some containers rename PID 1; skip rather than fail.
		t.Skipf("process check returned error (container environment?): %v", err)
	}
	if drifted {
		t.Errorf("expected no drift for process %q, got message: %s", processName, msg)
	}
}

func TestCheckProcessRunning_Drift(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping pgrep-based test on Windows")
	}

	// Use a process name that should never exist.
	drifted, msg, err := checkProcessRunning(map[string]string{
		"process": "driftwatch-nonexistent-proc-xyz",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Error("expected drift for non-existent process, got none")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}
