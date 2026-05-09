package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func rdsServer(status string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": status})
	}))
}

func TestCheckRDSInstanceStatus_MissingEndpoint(t *testing.T) {
	_, _, err := checkRDSInstanceStatus(map[string]string{
		"db_instance_identifier": "mydb",
	})
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckRDSInstanceStatus_MissingDBIdentifier(t *testing.T) {
	_, _, err := checkRDSInstanceStatus(map[string]string{
		"endpoint": "http://localhost:4566",
	})
	if err == nil {
		t.Fatal("expected error for missing db_instance_identifier")
	}
}

func TestCheckRDSInstanceStatus_NoDrift(t *testing.T) {
	srv := rdsServer("available")
	defer srv.Close()

	drift, _, err := checkRDSInstanceStatus(map[string]string{
		"endpoint":               srv.URL,
		"db_instance_identifier": "prod-db",
		"expected_status":        "available",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift")
	}
}

func TestCheckRDSInstanceStatus_Drift(t *testing.T) {
	srv := rdsServer("stopped")
	defer srv.Close()

	drift, msg, err := checkRDSInstanceStatus(map[string]string{
		"endpoint":               srv.URL,
		"db_instance_identifier": "prod-db",
		"expected_status":        "available",
	})
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

func TestCheckRDSInstanceStatus_DefaultExpectedAvailable(t *testing.T) {
	srv := rdsServer("available")
	defer srv.Close()

	drift, _, err := checkRDSInstanceStatus(map[string]string{
		"endpoint":               srv.URL,
		"db_instance_identifier": "prod-db",
		// no expected_status — should default to "available"
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift with default expected status")
	}
}

func TestCheckRDSInstanceStatus_ViaChecker(t *testing.T) {
	srv := rdsServer("available")
	defer srv.Close()

	c := New()
	drift, _, err := c.Check("rds_instance_status", map[string]string{
		"endpoint":               srv.URL,
		"db_instance_identifier": "replica-db",
		"expected_status":        "available",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift")
	}
}
