package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/example/driftwatch/internal/config"
)

// checkELBv2TargetGroupHealth checks the health of an ALB/NLB target group
// via the AWS ELBv2 API (or a compatible mock endpoint).
//
// Required fields:
//   - endpoint:         base URL of the ELBv2 API (e.g. http://localhost:4566)
//   - target_group_arn: ARN of the target group to inspect
//   - expected:         expected health state, e.g. "healthy" (default: "healthy")
func checkELBv2TargetGroupHealth(chk config.Check) (bool, string, error) {
	endpoint, ok := chk.Params["endpoint"].(string)
	if !ok || endpoint == "" {
		return false, "", fmt.Errorf("elbv2_target_group_health: missing 'endpoint'")
	}

	arn, ok := chk.Params["target_group_arn"].(string)
	if !ok || arn == "" {
		return false, "", fmt.Errorf("elbv2_target_group_health: missing 'target_group_arn'")
	}

	expected := "healthy"
	if v, ok := chk.Params["expected"].(string); ok && v != "" {
		expected = strings.ToLower(v)
	}

	url := fmt.Sprintf("%s/?Action=DescribeTargetHealth&TargetGroupArn=%s", endpoint, arn)
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return false, "", fmt.Errorf("elbv2_target_group_health: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("elbv2_target_group_health: read body: %w", err)
	}

	var result struct {
		State string `json:"State"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("elbv2_target_group_health: parse response: %w", err)
	}

	actual := strings.ToLower(result.State)
	if actual != expected {
		return true, fmt.Sprintf("target group %s state is %q, expected %q", arn, actual, expected), nil
	}
	return false, "", nil
}
