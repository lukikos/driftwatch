package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkEKSClusterStatus checks that an EKS cluster is in the expected status.
// Fields used from check.Fields:
//   - endpoint:       base URL for the EKS-compatible API (required)
//   - cluster_name:   name of the EKS cluster (required)
//   - expected:       expected cluster status (default: "ACTIVE")
func checkEKSClusterStatus(check config.Check) (bool, string, error) {
	endpoint, ok := check.Fields["endpoint"]
	if !ok || strings.TrimSpace(endpoint) == "" {
		return false, "", fmt.Errorf("eks_cluster_status: missing required field 'endpoint'")
	}

	clusterName, ok := check.Fields["cluster_name"]
	if !ok || strings.TrimSpace(clusterName) == "" {
		return false, "", fmt.Errorf("eks_cluster_status: missing required field 'cluster_name'")
	}

	expected := "ACTIVE"
	if v, ok := check.Fields["expected"]; ok && strings.TrimSpace(v) != "" {
		expected = strings.TrimSpace(v)
	}

	url := fmt.Sprintf("%s/clusters/%s", strings.TrimRight(endpoint, "/"), clusterName)
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return false, "", fmt.Errorf("eks_cluster_status: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("eks_cluster_status: failed to read response: %w", err)
	}

	var result struct {
		Cluster struct {
			Status string `json:"status"`
		} `json:"cluster"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("eks_cluster_status: failed to parse response: %w", err)
	}

	actual := result.Cluster.Status
	if actual != expected {
		return true, fmt.Sprintf("EKS cluster %q status is %q, expected %q", clusterName, actual, expected), nil
	}
	return false, "", nil
}
