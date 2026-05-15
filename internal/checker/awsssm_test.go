package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func ssmServer(value string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"value": value})
	}))
}

func TestCheckSSMParameter_MissingEndpoint(t *testing.T) {
	check := config.Check{Fields: map[string]interface{}{
		"parameter_name": "/app/env",
		"expected":       "production",
	}}
	_, _, err := checkSSMParameter(check)
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckSSMParameter_MissingParameterName(t *testing.T) {
	check := config.Check{Fields: map[string]interface{}{
		"endpoint": "http://localhost:4566",
		"expected": "production",
	}}
	_, _, err := checkSSMParameter(check)
	if err == nil {
		t.Fatal("expected error for missing parameter_name")
	}
}

func TestCheckSSMParameter_MissingExpected(t *testing.T) {
	check := config.Check{Fields: map[string]interface{}{
		"endpoint":       "http://localhost:4566",
		"parameter_name": "/app/env",
	}}
	_, _, err := checkSSMParameter(check)
	if err == nil {
		t.Fatal("expected error for missing expected")
	}
}

func TestCheckSSMParameter_NoDrift(t *testing.T) {
	srv := ssmServer("production", http.StatusOK)
	defer srv.Close()

	check := config.Check{Fields: map[string]interface{}{
		"endpoint":       srv.URL,
		"parameter_name": "/app/env",
		"expected":       "production",
	}}
	drift, _, err := checkSSMParameter(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift")
	}
}

func TestCheckSSMParameter_Drift(t *testing.T) {
	srv := ssmServer("staging", http.StatusOK)
	defer srv.Close()

	check := config.Check{Fields: map[string]interface{}{
		"endpoint":       srv.URL,
		"parameter_name": "/app/env",
		"expected":       "production",
	}}
	drift, msg, err := checkSSMParameter(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Fatal("expected drift")
	}
	if msg == "" {
		t.Fatal("expected non-empty drift message")
	}
}

func TestCheckSSMParameter_NotFound(t *testing.T) {
	srv := ssmServer("", http.StatusNotFound)
	defer srv.Close()

	check := config.Check{Fields: map[string]interface{}{
		"endpoint":       srv.URL,
		"parameter_name": "/app/missing",
		"expected":       "production",
	}}
	drift, msg, err := checkSSMParameter(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Fatal("expected drift when parameter not found")
	}
	if msg == "" {
		t.Fatal("expected non-empty drift message")
	}
}

func TestCheckSSMParameter_ViaChecker(t *testing.T) {
	srv := ssmServer("production", http.StatusOK)
	defer srv.Close()

	c := New()
	check := config.Check{
		Name: "ssm-env",
		Type: "ssm_parameter",
		Fields: map[string]interface{}{
			"endpoint":       srv.URL,
			"parameter_name": "/app/env",
			"expected":       "production",
		},
	}
	drift, _, err := c.Check(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift")
	}
}
