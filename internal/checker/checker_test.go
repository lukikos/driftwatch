package checker

import (
	"fmt"
	"net"
	"testing"

	"github.com/user/driftwatch/internal/config"
)

func TestCheckEnvVar_NoDrift(t *testing.T) {
	t.Setenv("TEST_VAR", "hello")
	c := New([]config.Check{{Name: "env", Type: "env_var", Fields: map[string]string{"name": "TEST_VAR", "expected": "hello"}}})
	results := c.RunAll()
	if results[0].Drifted {
		t.Errorf("expected no drift")
	}
}

func TestCheckEnvVar_Drift(t *testing.T) {
	t.Setenv("TEST_VAR", "world")
	c := New([]config.Check{{Name: "env", Type: "env_var", Fields: map[string]string{"name": "TEST_VAR", "expected": "hello"}}})
	results := c.RunAll()
	if !results[0].Drifted {
		t.Errorf("expected drift")
	}
}

func TestCheckFileHash_NoDrift(t *testing.T) {
	f := writeTempFile(t, "content")
	hash := sha256Hex("content")
	c := New([]config.Check{{Name: "file", Type: "file_hash", Fields: map[string]string{"path": f, "expected": hash}}})
	results := c.RunAll()
	if results[0].Drifted {
		t.Errorf("expected no drift")
	}
}

func TestCheckFileHash_Drift(t *testing.T) {
	f := writeTempFile(t, "changed")
	c := New([]config.Check{{Name: "file", Type: "file_hash", Fields: map[string]string{"path": f, "expected": "deadbeef"}}})
	results := c.RunAll()
	if !results[0].Drifted {
		t.Errorf("expected drift")
	}
}

func TestUnknownCheckType(t *testing.T) {
	c := New([]config.Check{{Name: "bad", Type: "nonexistent", Fields: map[string]string{}}})
	results := c.RunAll()
	if !results[0].Drifted {
		t.Errorf("expected drift for unknown type")
	}
}

func TestCheckPortOpen_ViaChecker_NoDrift(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()
	port := fmt.Sprintf("%d", ln.Addr().(*net.TCPAddr).Port)

	c := New([]config.Check{{
		Name:   "port",
		Type:   "port_open",
		Fields: map[string]string{"host": "127.0.0.1", "port": port},
	}})
	results := c.RunAll()
	if results[0].Drifted {
		t.Errorf("expected no drift for open port, got: %s", results[0].Message)
	}
}
