package checker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkRedshiftClusterStatus checks that an Amazon Redshift cluster is in the
// expected status (default: "available") by querying a mock-friendly endpoint
// that mirrors the AWS Redshift DescribeClusters API shape.
//
// Required fields:
//   - endpoint:        base URL for the Redshift API (e.g. http://localhost:4566)
//   - cluster_id:      Redshift cluster identifier
//
// Optional fields:
//   - expected_status: desired cluster status (default: "available")
func checkRedshiftClusterStatus(chk config.Check) (bool, string, error) {
	endpoint, ok := chk.Params["endpoint"]
	if !ok || strings.TrimSpace(endpoint) == "" {
		return false, "", fmt.Errorf("redshift_cluster_status: missing required param 'endpoint'")
	}

	clusterID, ok := chk.Params["cluster_id"]
	if !ok || strings.TrimSpace(clusterID) == "" {
		return false, "", fmt.Errorf("redshift_cluster_status: missing required param 'cluster_id'")
	}

	expected := "available"
	if v, ok := chk.Params["expected_status"]; ok && strings.TrimSpace(v) != "" {
		expected = strings.TrimSpace(v)
	}

	url := fmt.Sprintf("%s/redshift/clusters/%s", strings.TrimRight(endpoint, "/"), clusterID)
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return false, "", fmt.Errorf("redshift_cluster_status: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return true, fmt.Sprintf("cluster %q not found (drift: expected status %q)", clusterID, expected), nil
	}
	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("redshift_cluster_status: unexpected HTTP %d", resp.StatusCode)
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return false, "", fmt.Errorf("redshift_cluster_status: failed to decode response: %w", err)
	}

	if body.Status != expected {
		return true, fmt.Sprintf("cluster %q status is %q, expected %q", clusterID, body.Status, expected), nil
	}
	return false, fmt.Sprintf("cluster %q status is %q (ok)", clusterID, body.Status), nil
}
