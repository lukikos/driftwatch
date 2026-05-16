package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/yourusername/driftwatch/internal/config"
)

// checkElastiCacheClusterStatus checks whether an ElastiCache cluster
// reports the expected status via a mock-friendly HTTP endpoint.
//
// Required fields:
//   - endpoint: base URL of the ElastiCache-compatible API
//   - cluster_id: the ElastiCache cluster identifier
//
// Optional fields:
//   - expected: expected cluster status (default: "available")
func checkElastiCacheClusterStatus(c config.Check) (bool, string, error) {
	endpoint, ok := c.Params["endpoint"]
	if !ok || strings.TrimSpace(endpoint) == "" {
		return false, "", fmt.Errorf("elasticache_cluster: missing required param 'endpoint'")
	}

	clusterID, ok := c.Params["cluster_id"]
	if !ok || strings.TrimSpace(clusterID) == "" {
		return false, "", fmt.Errorf("elasticache_cluster: missing required param 'cluster_id'")
	}

	expected := "available"
	if v, ok := c.Params["expected"]; ok && strings.TrimSpace(v) != "" {
		expected = strings.TrimSpace(v)
	}

	url := fmt.Sprintf("%s/elasticache/clusters/%s", strings.TrimRight(endpoint, "/"), clusterID)
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return false, "", fmt.Errorf("elasticache_cluster: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("elasticache_cluster: failed to read response: %w", err)
	}

	var result struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("elasticache_cluster: failed to parse response: %w", err)
	}

	if result.Status != expected {
		return true, fmt.Sprintf("cluster %q status is %q, expected %q", clusterID, result.Status, expected), nil
	}
	return false, "", nil
}
