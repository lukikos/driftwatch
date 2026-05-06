package checker

import (
	"fmt"
	"net"
	"testing"
)

// startTCPListener binds to an ephemeral port and returns the port string and a closer.
func startTCPListener(t *testing.T) (string, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := fmt.Sprintf("%d", ln.Addr().(*net.TCPAddr).Port)
	return port, func() { ln.Close() }
}

func TestCheckPortOpen_NoDrift(t *testing.T) {
	port, close := startTCPListener(t)
	defer close()

	drift, msg, err := checkPortOpen(map[string]string{
		"host": "127.0.0.1",
		"port": port,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift, got message: %s", msg)
	}
}

func TestCheckPortOpen_Drift(t *testing.T) {
	// Use a port that is very unlikely to be open.
	drift, msg, err := checkPortOpen(map[string]string{
		"host": "127.0.0.1",
		"port": "19999",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Error("expected drift for closed port")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckPortOpen_DefaultHost(t *testing.T) {
	// No host provided — should default to localhost without error.
	// We don't assert drift state since port 19998 may or may not be open.
	_, _, err := checkPortOpen(map[string]string{
		"port": "19998",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckPortOpen_MissingPort(t *testing.T) {
	_, _, err := checkPortOpen(map[string]string{})
	if err == nil {
		t.Error("expected error for missing port field")
	}
}

func TestCheckPortOpen_EmptyPort(t *testing.T) {
	_, _, err := checkPortOpen(map[string]string{"port": ""})
	if err == nil {
		t.Error("expected error for empty port field")
	}
}
