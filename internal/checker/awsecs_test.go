package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/driftwatch/internal/config"
)

func ecsServer(status string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var resp struct {
			Services []struct {
				Status string `json:"status"`
			} `json:"services"`
		}
		if status != "" {
			resp.Services = []struct {
				Status string `json:"status"`
			}{{Status: status}}
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

func TestCheckECSServiceStatus_MissingEndpoint(t *testing.T) {
	_, _, err := checkECSServiceStatus(config.Check{Params: map[string]interface{}{}})
	if err == nil || err.Error() != "ecs_service_status: missing required param 'endpoint'" {
		t.Fatalf("expected endpoint error, got %v", err)
	}
}

func TestCheckECSServiceStatus_MissingCluster(t *testing.T) {
	_, _, err := checkECSServiceStatus(config.Check{Params: map[string]interface{}{"endpoint": "http://x"}})
	if err == nil || err.Error() != "ecs_service_status: missing required param 'cluster'" {
		t.Fatalf("expected cluster error, got %v", err)
	}
}

func TestCheckECSServiceStatus_MissingServiceName(t *testing.T) {
	_, _, err := checkECSServiceStatus(config.Check{Params: map[string]interface{}{
		"endpoint": "http://x",
		"cluster":  "prod",
	}})
	if err == nil || err.Error() != "ecs_service_status: missing required param 'service_name'" {
		t.Fatalf("expected service_name error, got %v", err)
	}
}

func TestCheckECSServiceStatus_NoDrift(t *testing.T) {
	srv := ecsServer("ACTIVE")
	defer srv.Close()

	drift, msg, err := checkECSServiceStatus(config.Check{Params: map[string]interface{}{
		"endpoint":     srv.URL,
		"cluster":      "prod",
		"service_name": "my-service",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatalf("expected no drift, got: %s", msg)
	}
}

func TestCheckECSServiceStatus_Drift(t *testing.T) {
	srv := ecsServer("DRAINING")
	defer srv.Close()

	drift, msg, err := checkECSServiceStatus(config.Check{Params: map[string]interface{}{
		"endpoint":        srv.URL,
		"cluster":         "prod",
		"service_name":    "my-service",
		"expected_status": "ACTIVE",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Fatal("expected drift but got none")
	}
	if msg == "" {
		t.Fatal("expected non-empty drift message")
	}
}

func TestCheckECSServiceStatus_DefaultExpectedActive(t *testing.T) {
	srv := ecsServer("ACTIVE")
	defer srv.Close()

	drift, _, err := checkECSServiceStatus(config.Check{Params: map[string]interface{}{
		"endpoint":     srv.URL,
		"cluster":      "staging",
		"service_name": "worker",
		// no expected_status — should default to ACTIVE
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift with default ACTIVE status")
	}
}

func TestCheckECSServiceStatus_NoServicesReturnedReportsDrift(t *testing.T) {
	srv := ecsServer("") // returns empty services list
	defer srv.Close()

	drift, msg, err := checkECSServiceStatus(config.Check{Params: map[string]interface{}{
		"endpoint":     srv.URL,
		"cluster":      "prod",
		"service_name": "missing-service",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Fatal("expected drift when service not found")
	}
	if msg == "" {
		t.Fatal("expected non-empty drift message")
	}
}
