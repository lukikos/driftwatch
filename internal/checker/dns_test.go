package checker

import (
	"testing"

	"github.com/user/driftwatch/internal/config"
)

func TestCheckDNSResolve_MissingHost(t *testing.T) {
	check := config.Check{
		Name: "dns-test",
		Type: "dns",
		Config: map[string]string{
			"expected": "127.0.0.1",
		},
	}
	_, _, err := checkDNSResolve(check)
	if err == nil {
		t.Fatal("expected error for missing host, got nil")
	}
}

func TestCheckDNSResolve_MissingExpected(t *testing.T) {
	check := config.Check{
		Name: "dns-test",
		Type: "dns",
		Config: map[string]string{
			"host": "localhost",
		},
	}
	_, _, err := checkDNSResolve(check)
	if err == nil {
		t.Fatal("expected error for missing expected, got nil")
	}
}

func TestCheckDNSResolve_NoDrift(t *testing.T) {
	// localhost should always resolve to 127.0.0.1
	check := config.Check{
		Name: "dns-localhost",
		Type: "dns",
		Config: map[string]string{
			"host":     "localhost",
			"expected": "127.0.0.1, ::1",
		},
	}
	drifted, _, err := checkDNSResolve(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Error("expected no drift for localhost resolving to 127.0.0.1 or ::1")
	}
}

func TestCheckDNSResolve_Drift(t *testing.T) {
	check := config.Check{
		Name: "dns-drift",
		Type: "dns",
		Config: map[string]string{
			"host":     "localhost",
			"expected": "1.2.3.4",
		},
	}
	drifted, msg, err := checkDNSResolve(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Error("expected drift when expected IP does not match")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckDNSResolve_ViaChecker(t *testing.T) {
	c := New()
	check := config.Check{
		Name: "dns-via-checker",
		Type: "dns",
		Config: map[string]string{
			"host":     "localhost",
			"expected": "127.0.0.1, ::1",
		},
	}
	drifted, _, err := c.Check(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Error("expected no drift")
	}
}
