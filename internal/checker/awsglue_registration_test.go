package checker

import (
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func TestGlueJobStatus_RegistrationInDispatch(t *testing.T) {
	c := New()
	// Providing an invalid endpoint ensures we reach the dispatch but fail
	// at the HTTP call, not at "unknown type".
	_, _, err := c.Check(config.Check{
		Name: "glue-reg-test",
		Type: "glue_job_status",
		Params: map[string]interface{}{
			"endpoint": "http://127.0.0.1:0", // unreachable
			"job_name": "some-job",
		},
	})
	if err == nil {
		t.Fatal("expected connection error, got nil")
	}
	if strings.Contains(err.Error(), "unknown check type") {
		t.Errorf("glue_job_status not registered in dispatch: %v", err)
	}
}

func TestGlueJobStatus_UnknownTypeStillErrors(t *testing.T) {
	c := New()
	_, _, err := c.Check(config.Check{
		Name: "unknown",
		Type: "glue_nonexistent_type",
		Params: map[string]interface{}{},
	})
	if err == nil {
		t.Fatal("expected error for unknown check type")
	}
	if !strings.Contains(err.Error(), "unknown check type") {
		t.Errorf("unexpected error message: %v", err)
	}
}
