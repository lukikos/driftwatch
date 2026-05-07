package checker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckHTTPLatency_MissingURL(t *testing.T) {
	_, _, err := checkHTTPLatency(map[string]string{"max_ms": "500"})
	if err == nil {
		t.Fatal("expected error for missing url")
	}
}

func TestCheckHTTPLatency_MissingMaxMS(t *testing.T) {
	_, _, err := checkHTTPLatency(map[string]string{"url": "http://example.com"})
	if err == nil {
		t.Fatal("expected error for missing max_ms")
	}
}

func TestCheckHTTPLatency_InvalidMaxMS(t *testing.T) {
	_, _, err := checkHTTPLatency(map[string]string{"url": "http://example.com", "max_ms": "notanumber"})
	if err == nil {
		t.Fatal("expected error for invalid max_ms")
	}
}

func TestCheckHTTPLatency_ZeroMaxMS(t *testing.T) {
	_, _, err := checkHTTPLatency(map[string]string{"url": "http://example.com", "max_ms": "0"})
	if err == nil {
		t.Fatal("expected error for zero max_ms")
	}
}

func TestCheckHTTPLatency_NoDrift(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	drift, msg, err := checkHTTPLatency(map[string]string{
		"url":    ts.URL,
		"max_ms": "5000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift for fast server, got: %s", msg)
	}
}

func TestCheckHTTPLatency_Drift(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(60 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	drift, msg, err := checkHTTPLatency(map[string]string{
		"url":    ts.URL,
		"max_ms": "10",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Errorf("expected drift for slow server, got: %s", msg)
	}
}

func TestCheckHTTPLatency_ViaChecker(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := New()
	drift, _, err := c.Check("http_latency", map[string]string{
		"url":    ts.URL,
		"max_ms": "5000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift")
	}
}

func TestCheckHTTPLatency_InvalidURL(t *testing.T) {
	_, _, err := checkHTTPLatency(map[string]string{
		"url":    "://bad-url",
		"max_ms": "500",
	})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
