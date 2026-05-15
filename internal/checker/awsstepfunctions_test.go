package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func stepFunctionsServer(status string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": status})
	}))
}

func TestCheckStepFunctionsStateMachine_MissingEndpoint(t *testing.T) {
	c := config.Check{Params: map[string]interface{}{"state_machine_arn": "arn:aws:states:us-east-1:123:stateMachine:MyMachine"}}
	_, _, err := checkStepFunctionsStateMachine(c)
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckStepFunctionsStateMachine_MissingARN(t *testing.T) {
	c := config.Check{Params: map[string]interface{}{"endpoint": "http://localhost"}}
	_, _, err := checkStepFunctionsStateMachine(c)
	if err == nil {
		t.Fatal("expected error for missing state_machine_arn")
	}
}

func TestCheckStepFunctionsStateMachine_NoDrift(t *testing.T) {
	srv := stepFunctionsServer("ACTIVE")
	defer srv.Close()

	c := config.Check{Params: map[string]interface{}{
		"endpoint":          srv.URL,
		"state_machine_arn": "arn:aws:states:us-east-1:123:stateMachine:MyMachine",
	}}
	drift, _, err := checkStepFunctionsStateMachine(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift for ACTIVE state machine")
	}
}

func TestCheckStepFunctionsStateMachine_Drift(t *testing.T) {
	srv := stepFunctionsServer("DELETING")
	defer srv.Close()

	c := config.Check{Params: map[string]interface{}{
		"endpoint":          srv.URL,
		"state_machine_arn": "arn:aws:states:us-east-1:123:stateMachine:MyMachine",
	}}
	drift, msg, err := checkStepFunctionsStateMachine(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Error("expected drift for DELETING state machine")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckStepFunctionsStateMachine_DefaultExpectedActive(t *testing.T) {
	srv := stepFunctionsServer("ACTIVE")
	defer srv.Close()

	c := config.Check{Params: map[string]interface{}{
		"endpoint":          srv.URL,
		"state_machine_arn": "arn:aws:states:us-east-1:123:stateMachine:MyMachine",
		// no expected_status — should default to ACTIVE
	}}
	drift, _, err := checkStepFunctionsStateMachine(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift when default expected_status matches")
	}
}

func TestCheckStepFunctionsStateMachine_ViaChecker(t *testing.T) {
	srv := stepFunctionsServer("ACTIVE")
	defer srv.Close()

	ch := New()
	c := config.Check{
		Name: "sfn-test",
		Type: "step_functions_state_machine",
		Params: map[string]interface{}{
			"endpoint":          srv.URL,
			"state_machine_arn": "arn:aws:states:us-east-1:123:stateMachine:MyMachine",
		},
	}
	drift, _, err := ch.Run(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift")
	}
}
