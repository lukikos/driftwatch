package checker

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkRoute53HealthCheck queries an AWS Route 53 health check endpoint
// and compares the returned status against the expected value.
//
// Required fields:
//   - endpoint: base URL for the Route 53 API (e.g. https://route53.amazonaws.com)
//   - health_check_id: the Route 53 health check ID
//   - expected: expected status string (default: "Success")
func checkRoute53HealthCheck(chk config.Check) (bool, string, error) {
	endpoint, ok := chk.Fields["endpoint"]
	if !ok || strings.TrimSpace(endpoint) == "" {
		return false, "", fmt.Errorf("route53_health_check: missing required field 'endpoint'")
	}

	healthCheckID, ok := chk.Fields["health_check_id"]
	if !ok || strings.TrimSpace(healthCheckID) == "" {
		return false, "", fmt.Errorf("route53_health_check: missing required field 'health_check_id'")
	}

	expected := chk.Fields["expected"]
	if strings.TrimSpace(expected) == "" {
		expected = "Success"
	}

	url := fmt.Sprintf("%s/2013-04-01/healthcheck/%s/status",
		strings.TrimRight(endpoint, "/"), healthCheckID)

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return false, "", fmt.Errorf("route53_health_check: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("route53_health_check: failed to read response: %w", err)
	}

	var result struct {
		CheckerReports []struct {
			Status string `xml:"Status"`
		} `xml:"CheckerReport"`
	}

	if err := xml.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("route53_health_check: failed to parse XML: %w", err)
	}

	if len(result.CheckerReports) == 0 {
		return true, fmt.Sprintf("no checker reports found for health check %s", healthCheckID), nil
	}

	actual := strings.TrimSpace(result.CheckerReports[0].Status)
	if !strings.EqualFold(actual, expected) {
		return true, fmt.Sprintf("route53 health check %s: expected status %q, got %q",
			healthCheckID, expected, actual), nil
	}

	return false, "", nil
}
