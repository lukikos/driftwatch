package checker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHTTPLatency_RegistrationInDispatch ensures that the http_latency type is
// wired into the checker dispatch table and returns a sensible result.
func TestHTTPLatency_RegistrationInDispatch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := New()

	drift, msg, err := c.Check("http_latency", map[string]string{
		"url":    ts.URL,
		"max_ms": "9999",
	})
	if err != nil {
		t.Fatalf("http_latency check returned unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift for generous threshold, got message: %s", msg)
	}
}

// TestHTTPLatency_UnknownTypeStillErrors verifies that an unregistered type
// still returns an error (regression guard for the dispatch table).
func TestHTTPLatency_UnknownTypeStillErrors(t *testing.T) {
	c := New()
	_, _, err := c.Check("nonexistent_type_xyz", map[string]string{})
	if err == nil {
		t.Fatal("expected error for unknown check type")
	}
}
