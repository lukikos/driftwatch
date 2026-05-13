package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/user/driftwatch/internal/config"
)

// checkECSServiceStatus checks whether an ECS service's status matches the expected value.
// It queries a local or mock ECS-compatible endpoint (or AWS ECS DescribeServices API).
// Required fields: endpoint, cluster, service_name, expected_status (default: "ACTIVE")
func checkECSServiceStatus(check config.Check) (bool, string, error) {
	endpoint, _ := check.Params["endpoint"].(string)
	cluster, _ := check.Params["cluster"].(string)
	serviceName, _ := check.Params["service_name"].(string)
	expectedStatus, _ := check.Params["expected_status"].(string)

	if endpoint == "" {
		return false, "", fmt.Errorf("ecs_service_status: missing required param 'endpoint'")
	}
	if cluster == "" {
		return false, "", fmt.Errorf("ecs_service_status: missing required param 'cluster'")
	}
	if serviceName == "" {
		return false, "", fmt.Errorf("ecs_service_status: missing required param 'service_name'")
	}
	if expectedStatus == "" {
		expectedStatus = "ACTIVE"
	}

	url := strings.TrimRight(endpoint, "/") +
		fmt.Sprintf("/?Action=DescribeServices&cluster=%s&services=%s", cluster, serviceName)

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return false, "", fmt.Errorf("ecs_service_status: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("ecs_service_status: failed to read response: %w", err)
	}

	var result struct {
		Services []struct {
			Status string `json:"status"`
		} `json:"services"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("ecs_service_status: failed to parse response: %w", err)
	}

	if len(result.Services) == 0 {
		return true, fmt.Sprintf("no ECS service '%s' found in cluster '%s'", serviceName, cluster), nil
	}

	actual := result.Services[0].Status
	if actual != expectedStatus {
		return true, fmt.Sprintf("ECS service '%s' status is '%s', expected '%s'", serviceName, actual, expectedStatus), nil
	}

	return false, "", nil
}
