package checker

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func makeCmd(name, command, expected string) config.Check {
	return config.Check{
		Name: name,
		Type: "syscmd",
		Config: map[string]string{
			"command":  command,
			"expected": expected,
		},
	}
}

func TestCheckSysCommand_NoDrift(t *testing.T) {
	check := makeCmd("echo-hello", "echo hello", "hello")
	drift, msg, err := checkSysCommand(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift, got message: %s", msg)
	}
}

func TestCheckSysCommand_Drift(t *testing.T) {
	check := makeCmd("echo-mismatch", "echo world", "hello")
	drift, msg, err := checkSysCommand(check)
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

func TestCheckSysCommand_MissingCommand(t *testing.T) {
	check := config.Check{
		Name: "no-cmd",
		Type: "syscmd",
		Config: map[string]string{"expected": "foo"},
	}
	_, _, err := checkSysCommand(check)
	if err == nil {
		t.Error("expected error for missing command")
	}
}

func TestCheckSysCommand_MissingExpected(t *testing.T) {
	check := config.Check{
		Name: "no-expected",
		Type: "syscmd",
		Config: map[string]string{"command": "echo hi"},
	}
	_, _, err := checkSysCommand(check)
	if err == nil {
		t.Error("expected error for missing expected field")
	}
}

func TestCheckSysCommand_CommandFails(t *testing.T) {
	check := makeCmd("bad-cmd", "exit 1", "")
	_, _, err := checkSysCommand(check)
	if err == nil {
		t.Error("expected error when command exits non-zero")
	}
}

func TestCheckSysCommand_ViaChecker(t *testing.T) {
	c := New()
	check := makeCmd("via-checker", "echo driftwatch", "driftwatch")
	drift, msg, err := c.Check(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift, got: %s", msg)
	}
}
