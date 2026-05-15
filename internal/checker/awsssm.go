package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkSSMParameter fetches an AWS SSM Parameter Store value via a mock-friendly
// HTTP endpoint and compares it against the expected value.
//
// Required fields:
//   - endpoint: base URL of the SSM-compatible API (e.g. http://localhost:4566)
//   - parameter_name: the SSM parameter name (e.g. /app/env)
//   - expected: the expected string value of the parameter
func checkSSMParameter(check config.Check) (bool, string, error) {
	endpoint, ok := check.Fields["endpoint"].(string)
	if !ok || endpoint == "" {
		return false, "", fmt.Errorf("ssm_parameter: missing or empty 'endpoint'")
	}

	paramName, ok := check.Fields["parameter_name"].(string)
	if !ok || paramName == "" {
		return false, "", fmt.Errorf("ssm_parameter: missing or empty 'parameter_name'")
	}

	expected, ok := check.Fields["expected"].(string)
	if !ok || expected == "" {
		return false, "", fmt.Errorf("ssm_parameter: missing or empty 'expected'")
	}

	url := fmt.Sprintf("%s/ssm/parameters/%s", endpoint, paramName)
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return false, "", fmt.Errorf("ssm_parameter: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return true, fmt.Sprintf("parameter %q not found", paramName), nil
	}
	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("ssm_parameter: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("ssm_parameter: failed to read body: %w", err)
	}

	var result struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("ssm_parameter: failed to parse response: %w", err)
	}

	if result.Value != expected {
		return true, fmt.Sprintf("parameter %q value %q does not match expected %q", paramName, result.Value, expected), nil
	}

	return false, "", nil
}
