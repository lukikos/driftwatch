package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkStepFunctionsStateMachine checks the status of an AWS Step Functions
// state machine via the provided endpoint (real or mock).
//
// Required fields:
//   - endpoint:       base URL of the Step Functions API (or mock)
//   - state_machine_arn: ARN of the state machine to inspect
//
// Optional fields:
//   - expected_status: desired status string (default: "ACTIVE")
func checkStepFunctionsStateMachine(c config.Check) (bool, string, error) {
	endpoint, ok := c.Params["endpoint"].(string)
	if !ok || strings.TrimSpace(endpoint) == "" {
		return false, "", fmt.Errorf("step_functions_state_machine: missing or empty 'endpoint'")
	}

	arn, ok := c.Params["state_machine_arn"].(string)
	if !ok || strings.TrimSpace(arn) == "" {
		return false, "", fmt.Errorf("step_functions_state_machine: missing or empty 'state_machine_arn'")
	}

	expected := "ACTIVE"
	if v, ok := c.Params["expected_status"].(string); ok && strings.TrimSpace(v) != "" {
		expected = strings.TrimSpace(v)
	}

	url := fmt.Sprintf("%s/statemachines/%s", strings.TrimRight(endpoint, "/"), arn)
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return false, "", fmt.Errorf("step_functions_state_machine: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("step_functions_state_machine: failed to read response: %w", err)
	}

	var result struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("step_functions_state_machine: failed to parse response: %w", err)
	}

	if result.Status != expected {
		return true, fmt.Sprintf("state machine %q status is %q, expected %q", arn, result.Status, expected), nil
	}
	return false, "", nil
}
