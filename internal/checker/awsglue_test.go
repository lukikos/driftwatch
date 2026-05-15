package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func glueServer(t *testing.T, jobName, status string, statusCode int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if statusCode == http.StatusOK {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Job":       map[string]interface{}{"Name": jobName},
				"JobStatus": status,
			})
		}
	}))
}

func TestCheckGlueJobStatus_MissingEndpoint(t *testing.T) {
	_, _, err := checkGlueJobStatus(config.Check{Params: map[string]interface{}{"job_name": "myjob"}})
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckGlueJobStatus_MissingJobName(t *testing.T) {
	_, _, err := checkGlueJobStatus(config.Check{Params: map[string]interface{}{"endpoint": "http://localhost"}})
	if err == nil {
		t.Fatal("expected error for missing job_name")
	}
}

func TestCheckGlueJobStatus_NoDrift(t *testing.T) {
	srv := glueServer(t, "etl-job", "READY", http.StatusOK)
	defer srv.Close()

	drift, msg, err := checkGlueJobStatus(config.Check{Params: map[string]interface{}{
		"endpoint": srv.URL,
		"job_name": "etl-job",
		"expected": "READY",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift, got message: %s", msg)
	}
}

func TestCheckGlueJobStatus_Drift(t *testing.T) {
	srv := glueServer(t, "etl-job", "FAILED", http.StatusOK)
	defer srv.Close()

	drift, msg, err := checkGlueJobStatus(config.Check{Params: map[string]interface{}{
		"endpoint": srv.URL,
		"job_name": "etl-job",
		"expected": "READY",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Error("expected drift but got none")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckGlueJobStatus_JobNotFound(t *testing.T) {
	srv := glueServer(t, "", "", http.StatusNotFound)
	defer srv.Close()

	drift, msg, err := checkGlueJobStatus(config.Check{Params: map[string]interface{}{
		"endpoint": srv.URL,
		"job_name": "missing-job",
		"expected": "READY",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Error("expected drift for missing job")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckGlueJobStatus_DefaultExpectedREADY(t *testing.T) {
	srv := glueServer(t, "etl-job", "READY", http.StatusOK)
	defer srv.Close()

	drift, _, err := checkGlueJobStatus(config.Check{Params: map[string]interface{}{
		"endpoint": srv.URL,
		"job_name": "etl-job",
		// no expected — should default to READY
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift with default expected=READY")
	}
}

func TestCheckGlueJobStatus_ViaChecker(t *testing.T) {
	srv := glueServer(t, "my-job", "READY", http.StatusOK)
	defer srv.Close()

	c := New()
	drift, _, err := c.Check(config.Check{
		Name: "glue-test",
		Type: "glue_job_status",
		Params: map[string]interface{}{
			"endpoint": srv.URL,
			"job_name": "my-job",
			"expected": "READY",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift")
	}
}
