package checker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckHTTPStatus_NoDrift(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	drifted, msg, err := checkHTTPStatus(ts.URL, http.StatusOK)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Errorf("expected no drift, got message: %s", msg)
	}
}

func TestCheckHTTPStatus_Drift(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	drifted, msg, err := checkHTTPStatus(ts.URL, http.StatusOK)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Error("expected drift but got none")
	}
	if msg == "" {
		t.Error("expected a drift message")
	}
}

func TestCheckHTTPStatus_DefaultExpected200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// expectedStatus=0 should default to 200
	drifted, _, err := checkHTTPStatus(ts.URL, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Error("expected no drift with default expected status")
	}
}

func TestCheckHTTPStatus_EmptyURL(t *testing.T) {
	_, _, err := checkHTTPStatus("", http.StatusOK)
	if err == nil {
		t.Error("expected error for empty URL")
	}
}

func TestCheckHTTPStatus_InvalidURL(t *testing.T) {
	_, _, err := checkHTTPStatus("http://127.0.0.1:0/unreachable", http.StatusOK)
	if err == nil {
		t.Error("expected error for unreachable URL")
	}
}
