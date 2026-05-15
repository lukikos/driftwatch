package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func cloudFrontServer(distributionID, status string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, distributionID) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
			"Distribution": map[string]string{"Status": status},
		})
	}))
}

func TestCheckCloudFrontDistribution_MissingEndpoint(t *testing.T) {
	check := config.Check{Type: "cloudfront_distribution", Fields: map[string]string{
		"distribution_id": "E1234",
	}}
	_, _, err := checkCloudFrontDistribution(check)
	if err == nil || !strings.Contains(err.Error(), "endpoint") {
		t.Fatalf("expected endpoint error, got %v", err)
	}
}

func TestCheckCloudFrontDistribution_MissingDistributionID(t *testing.T) {
	check := config.Check{Type: "cloudfront_distribution", Fields: map[string]string{
		"endpoint": "http://localhost",
	}}
	_, _, err := checkCloudFrontDistribution(check)
	if err == nil || !strings.Contains(err.Error(), "distribution_id") {
		t.Fatalf("expected distribution_id error, got %v", err)
	}
}

func TestCheckCloudFrontDistribution_NoDrift(t *testing.T) {
	srv := cloudFrontServer("E1234", "Deployed")
	defer srv.Close()

	check := config.Check{Type: "cloudfront_distribution", Fields: map[string]string{
		"endpoint":        srv.URL,
		"distribution_id": "E1234",
	}}
	drift, msg, err := checkCloudFrontDistribution(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatalf("expected no drift, got: %s", msg)
	}
}

func TestCheckCloudFrontDistribution_Drift(t *testing.T) {
	srv := cloudFrontServer("E1234", "InProgress")
	defer srv.Close()

	check := config.Check{Type: "cloudfront_distribution", Fields: map[string]string{
		"endpoint":        srv.URL,
		"distribution_id": "E1234",
	}}
	drift, msg, err := checkCloudFrontDistribution(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Fatal("expected drift but got none")
	}
	if !strings.Contains(msg, "InProgress") {
		t.Fatalf("expected msg to contain status, got: %s", msg)
	}
}

func TestCheckCloudFrontDistribution_DefaultExpectedDeployed(t *testing.T) {
	srv := cloudFrontServer("EABC", "Deployed")
	defer srv.Close()

	check := config.Check{Type: "cloudfront_distribution", Fields: map[string]string{
		"endpoint":        srv.URL,
		"distribution_id": "EABC",
		// no 'expected' field — should default to "Deployed"
	}}
	drift, _, err := checkCloudFrontDistribution(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift with default expected=Deployed")
	}
}

func TestCheckCloudFrontDistribution_NotFound(t *testing.T) {
	srv := cloudFrontServer("E9999", "Deployed")
	defer srv.Close()

	check := config.Check{Type: "cloudfront_distribution", Fields: map[string]string{
		"endpoint":        srv.URL,
		"distribution_id": "EUNKNOWN",
	}}
	drift, msg, err := checkCloudFrontDistribution(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Fatal("expected drift for missing distribution")
	}
	if !strings.Contains(msg, "not found") {
		t.Fatalf("expected 'not found' in msg, got: %s", msg)
	}
}
