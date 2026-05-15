package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkCloudFrontDistribution checks whether a CloudFront distribution's status
// matches the expected value. By default it expects "Deployed".
//
// Required fields:
//   - endpoint: base URL for the CloudFront-compatible API (e.g. http://localhost:4566)
//   - distribution_id: the CloudFront distribution ID
//
// Optional fields:
//   - expected: expected status string (default: "Deployed")
func checkCloudFrontDistribution(check config.Check) (bool, string, error) {
	endpoint, ok := check.Fields["endpoint"]
	if !ok || strings.TrimSpace(endpoint) == "" {
		return false, "", fmt.Errorf("cloudfront_distribution: missing required field 'endpoint'")
	}

	distributionID, ok := check.Fields["distribution_id"]
	if !ok || strings.TrimSpace(distributionID) == "" {
		return false, "", fmt.Errorf("cloudfront_distribution: missing required field 'distribution_id'")
	}

	expected := "Deployed"
	if v, ok := check.Fields["expected"]; ok && strings.TrimSpace(v) != "" {
		expected = strings.TrimSpace(v)
	}

	url := fmt.Sprintf("%s/2020-05-31/distribution/%s", strings.TrimRight(endpoint, "/"), distributionID)

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return false, "", fmt.Errorf("cloudfront_distribution: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return true, fmt.Sprintf("distribution %q not found", distributionID), nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("cloudfront_distribution: failed to read response: %w", err)
	}

	var result struct {
		Distribution struct {
			Status string `json:"Status"`
		} `json:"Distribution"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("cloudfront_distribution: failed to parse response: %w", err)
	}

	actual := result.Distribution.Status
	if !strings.EqualFold(actual, expected) {
		return true, fmt.Sprintf("distribution %q status is %q, expected %q", distributionID, actual, expected), nil
	}

	return false, "", nil
}
