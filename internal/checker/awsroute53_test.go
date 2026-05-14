package checker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func route53Server(status string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<HealthCheckStatus>
  <CheckerReport>
    <Status>%s</Status>
  </CheckerReport>
</HealthCheckStatus>`, status)
	}))
}

func TestCheckRoute53HealthCheck_MissingEndpoint(t *testing.T) {
	chk := config.Check{Type: "route53_health_check", Fields: map[string]string{
		"health_check_id": "abc123",
	}}
	_, _, err := checkRoute53HealthCheck(chk)
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckRoute53HealthCheck_MissingHealthCheckID(t *testing.T) {
	chk := config.Check{Type: "route53_health_check", Fields: map[string]string{
		"endpoint": "http://localhost",
	}}
	_, _, err := checkRoute53HealthCheck(chk)
	if err == nil {
		t.Fatal("expected error for missing health_check_id")
	}
}

func TestCheckRoute53HealthCheck_NoDrift(t *testing.T) {
	ts := route53Server("Success")
	defer ts.Close()

	chk := config.Check{Type: "route53_health_check", Fields: map[string]string{
		"endpoint":        ts.URL,
		"health_check_id": "hc-001",
		"expected":        "Success",
	}}

	drifted, msg, err := checkRoute53HealthCheck(chk)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Errorf("expected no drift, got message: %s", msg)
	}
}

func TestCheckRoute53HealthCheck_Drift(t *testing.T) {
	ts := route53Server("Failure")
	defer ts.Close()

	chk := config.Check{Type: "route53_health_check", Fields: map[string]string{
		"endpoint":        ts.URL,
		"health_check_id": "hc-001",
		"expected":        "Success",
	}}

	drifted, msg, err := checkRoute53HealthCheck(chk)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Error("expected drift but got none")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckRoute53HealthCheck_DefaultExpectedSuccess(t *testing.T) {
	ts := route53Server("Success")
	defer ts.Close()

	chk := config.Check{Type: "route53_health_check", Fields: map[string]string{
		"endpoint":        ts.URL,
		"health_check_id": "hc-002",
		// no 'expected' field — should default to "Success"
	}}

	drifted, _, err := checkRoute53HealthCheck(chk)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Error("expected no drift with default expected=Success")
	}
}

func TestCheckRoute53HealthCheck_ViaChecker(t *testing.T) {
	ts := route53Server("Success")
	defer ts.Close()

	c := New()
	chk := config.Check{
		Name: "r53-test",
		Type: "route53_health_check",
		Fields: map[string]string{
			"endpoint":        ts.URL,
			"health_check_id": "hc-003",
		},
	}

	drifted, _, err := c.Check(chk)
	if err != nil {
		t.Fatalf("unexpected error via Checker: %v", err)
	}
	if drifted {
		t.Error("expected no drift via Checker")
	}
}
