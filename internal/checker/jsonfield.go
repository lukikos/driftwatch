package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// checkJSONField fetches a JSON endpoint and checks that a specific field
// matches an expected value. Supports dot-notation for nested fields.
//
// Required fields:
//   - url:      HTTP(S) URL to fetch
//   - field:    dot-separated path to the JSON field (e.g. "status" or "db.connected")
//   - expected: expected string value of the field
func checkJSONField(fields map[string]string) (bool, string, error) {
	url, ok := fields["url"]
	if !ok || url == "" {
		return false, "", fmt.Errorf("jsonfield check requires 'url'")
	}
	field, ok := fields["field"]
	if !ok || field == "" {
		return false, "", fmt.Errorf("jsonfield check requires 'field'")
	}
	expected, ok := fields["expected"]
	if !ok {
		return false, "", fmt.Errorf("jsonfield check requires 'expected'")
	}

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return false, "", fmt.Errorf("jsonfield: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("jsonfield: failed to read response: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return false, "", fmt.Errorf("jsonfield: failed to parse JSON: %w", err)
	}

	actual, err := extractField(data, strings.Split(field, "."))
	if err != nil {
		return false, "", fmt.Errorf("jsonfield: %w", err)
	}

	actualStr := fmt.Sprintf("%v", actual)
	if actualStr != expected {
		return true, fmt.Sprintf("field %q: got %q, expected %q", field, actualStr, expected), nil
	}
	return false, "", nil
}

func extractField(data map[string]interface{}, keys []string) (interface{}, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("empty field path")
	}
	val, ok := data[keys[0]]
	if !ok {
		return nil, fmt.Errorf("field %q not found in response", keys[0])
	}
	if len(keys) == 1 {
		return val, nil
	}
	nested, ok := val.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("field %q is not an object, cannot traverse deeper", keys[0])
	}
	return extractField(nested, keys[1:])
}
