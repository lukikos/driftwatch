package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func redshiftServer(status string, httpCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if httpCode != http.StatusOK {
			w.WriteHeader(httpCode)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": status})
	}))
}

func TestCheckRedshiftClusterStatus_MissingEndpoint(t *testing.T) {
	chk := config.Check{Type: "redshift_cluster_status", Params: map[string]string{"cluster_id": "my-cluster"}}
	_, _, err := checkRedshiftClusterStatus(chk)
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckRedshiftClusterStatus_MissingClusterID(t *testing.T) {
	chk := config.Check{Type: "redshift_cluster_status", Params: map[string]string{"endpoint": "http://localhost"}}
	_, _, err := checkRedshiftClusterStatus(chk)
	if err == nil {
		t.Fatal("expected error for missing cluster_id")
	}
}

func TestCheckRedshiftClusterStatus_NoDrift(t *testing.T) {
	srv := redshiftServer("available", http.StatusOK)
	defer srv.Close()

	chk := config.Check{
		Type: "redshift_cluster_status",
		Params: map[string]string{"endpoint": srv.URL, "cluster_id": "my-cluster"},
	}
	drifted, msg, err := checkRedshiftClusterStatus(chk)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Errorf("expected no drift, got: %s", msg)
	}
}

func TestCheckRedshiftClusterStatus_Drift(t *testing.T) {
	srv := redshiftServer("paused", http.StatusOK)
	defer srv.Close()

	chk := config.Check{
		Type: "redshift_cluster_status",
		Params: map[string]string{"endpoint": srv.URL, "cluster_id": "my-cluster"},
	}
	drifted, msg, err := checkRedshiftClusterStatus(chk)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Errorf("expected drift, got msg: %s", msg)
	}
}

func TestCheckRedshiftClusterStatus_DefaultExpectedAvailable(t *testing.T) {
	srv := redshiftServer("available", http.StatusOK)
	defer srv.Close()

	chk := config.Check{
		Type:   "redshift_cluster_status",
		Params: map[string]string{"endpoint": srv.URL, "cluster_id": "my-cluster"},
	}
	drifted, _, err := checkRedshiftClusterStatus(chk)
	if err != nil || drifted {
		t.Errorf("expected no drift with default expected_status=available")
	}
}

func TestCheckRedshiftClusterStatus_NotFound(t *testing.T) {
	srv := redshiftServer("", http.StatusNotFound)
	defer srv.Close()

	chk := config.Check{
		Type:   "redshift_cluster_status",
		Params: map[string]string{"endpoint": srv.URL, "cluster_id": "ghost-cluster"},
	}
	drifted, msg, err := checkRedshiftClusterStatus(chk)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Errorf("expected drift for not-found cluster, got: %s", msg)
	}
}
