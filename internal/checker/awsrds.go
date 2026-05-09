package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// checkRDSInstanceStatus checks an AWS RDS instance's status via a mock/proxy
// endpoint or the AWS RDS describe-db-instances API (localstack-compatible).
//
// Required fields:
//   - endpoint: base URL of the RDS-compatible API (e.g. http://localhost:4566)
//   - db_instance_identifier: the RDS instance identifier
//   - expected_status: expected DB status string (e.g. "available")
func checkRDSInstanceStatus(fields map[string]string) (bool, string, error) {
	endpoint, ok := fields["endpoint"]
	if !ok || strings.TrimSpace(endpoint) == "" {
		return false, "", fmt.Errorf("rds_instance_status: missing required field 'endpoint'")
	}

	dbID, ok := fields["db_instance_identifier"]
	if !ok || strings.TrimSpace(dbID) == "" {
		return false, "", fmt.Errorf("rds_instance_status: missing required field 'db_instance_identifier'")
	}

	expected := fields["expected_status"]
	if strings.TrimSpace(expected) == "" {
		expected = "available"
	}

	url := fmt.Sprintf("%s/rds/instances/%s", strings.TrimRight(endpoint, "/"), dbID)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false, "", fmt.Errorf("rds_instance_status: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("rds_instance_status: failed to read response: %w", err)
	}

	var result struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("rds_instance_status: failed to parse response: %w", err)
	}

	actual := strings.TrimSpace(result.Status)
	if !strings.EqualFold(actual, expected) {
		return true, fmt.Sprintf("RDS instance %q status is %q, expected %q", dbID, actual, expected), nil
	}

	return false, "", nil
}
