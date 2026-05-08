package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkLambdaFunction queries an AWS Lambda function's configuration via the
// AWS Lambda API (or a compatible mock endpoint) and checks whether the
// function's State matches the expected value.
//
// Required fields:
//   - function_name: the Lambda function name or ARN
//   - endpoint:      base URL of the Lambda API (e.g. http://localhost:3000)
//
// Optional fields:
//   - expected_state: expected function State (default: "Active")
func checkLambdaFunction(check config.Check) (bool, string, error) {
	functionName := strings.TrimSpace(check.Fields["function_name"])
	if functionName == "" {
		return false, "", fmt.Errorf("lambda_function: missing required field 'function_name'")
	}

	endpoint := strings.TrimSpace(check.Fields["endpoint"])
	if endpoint == "" {
		return false, "", fmt.Errorf("lambda_function: missing required field 'endpoint'")
	}

	expected := strings.TrimSpace(check.Fields["expected_state"])
	if expected == "" {
		expected = "Active"
	}

	url := fmt.Sprintf("%s/2015-03-31/functions/%s/configuration",
		strings.TrimRight(endpoint, "/"), functionName)

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return false, "", fmt.Errorf("lambda_function: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("lambda_function: failed to read response: %w", err)
	}

	var payload struct {
		State string `json:"State"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return false, "", fmt.Errorf("lambda_function: failed to parse response: %w", err)
	}

	actual := payload.State
	if actual != expected {
		return true, fmt.Sprintf("lambda function %q state is %q, expected %q", functionName, actual, expected), nil
	}
	return false, fmt.Sprintf("lambda function %q state is %q", functionName, actual), nil
}
