package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// checkDynamoDBTable queries a DynamoDB-compatible endpoint for a table's status
// and compares it against the expected value.
//
// Required fields:
//   - endpoint: base URL of the DynamoDB-compatible API
//   - table_name: name of the DynamoDB table
//   - expected: expected table status (e.g. "ACTIVE")
func checkDynamoDBTable(fields map[string]string) (bool, string, error) {
	endpoint, ok := fields["endpoint"]
	if !ok || strings.TrimSpace(endpoint) == "" {
		return false, "", fmt.Errorf("dynamodb_table: missing required field 'endpoint'")
	}

	tableName, ok := fields["table_name"]
	if !ok || strings.TrimSpace(tableName) == "" {
		return false, "", fmt.Errorf("dynamodb_table: missing required field 'table_name'")
	}

	expected, ok := fields["expected"]
	if !ok || strings.TrimSpace(expected) == "" {
		expected = "ACTIVE"
	}

	url := fmt.Sprintf("%s/tables/%s", strings.TrimRight(endpoint, "/"), tableName)
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return false, "", fmt.Errorf("dynamodb_table: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("dynamodb_table: failed to read response: %w", err)
	}

	var result struct {
		Table struct {
			TableStatus string `json:"TableStatus"`
		} `json:"Table"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("dynamodb_table: failed to parse response: %w", err)
	}

	actual := result.Table.TableStatus
	if actual != expected {
		return true, fmt.Sprintf("table %q status is %q, expected %q", tableName, actual, expected), nil
	}
	return false, "", nil
}
