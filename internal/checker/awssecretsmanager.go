package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/user/driftwatch/internal/config"
)

// checkSecretsManagerSecret verifies that an AWS Secrets Manager secret exists
// and optionally that a specific metadata field matches an expected value.
//
// Required fields:
//   - endpoint: base URL of the Secrets Manager API (or mock)
//   - secret_id: the name or ARN of the secret
//
// Optional fields:
//   - field:    a top-level JSON key from the DescribeSecret response to inspect
//   - expected: expected value for that field (default: "true" when field omitted)
func checkSecretsManagerSecret(check config.Check) (bool, string, error) {
	endpoint, ok := check.Params["endpoint"]
	if !ok || strings.TrimSpace(endpoint) == "" {
		return false, "", fmt.Errorf("secrets_manager_secret: missing required param 'endpoint'")
	}

	secretID, ok := check.Params["secret_id"]
	if !ok || strings.TrimSpace(secretID) == "" {
		return false, "", fmt.Errorf("secrets_manager_secret: missing required param 'secret_id'")
	}

	url := fmt.Sprintf("%s/secretsmanager/get-secret-value?secretId=%s", strings.TrimRight(endpoint, "/"), secretID)
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return false, "", fmt.Errorf("secrets_manager_secret: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return true, fmt.Sprintf("secret %q not found (404)", secretID), nil
	}
	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("secrets_manager_secret: unexpected status %d", resp.StatusCode)
	}

	field, hasField := check.Params["field"]
	if !hasField {
		// Secret exists — no drift.
		return false, "", nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("secrets_manager_secret: failed to read response: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return false, "", fmt.Errorf("secrets_manager_secret: failed to parse response JSON: %w", err)
	}

	actual, exists := data[field]
	if !exists {
		return false, "", fmt.Errorf("secrets_manager_secret: field %q not found in response", field)
	}

	expected := check.Params["expected"]
	actualStr := fmt.Sprintf("%v", actual)
	if actualStr != expected {
		return true, fmt.Sprintf("secret %q field %q: got %q, want %q", secretID, field, actualStr, expected), nil
	}

	return false, "", nil
}
