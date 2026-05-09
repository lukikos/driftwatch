package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// checkIAMRolePolicy checks that an AWS IAM role has an expected policy attached.
// It uses a configurable endpoint (for local/mock testing) or the real AWS IAM API.
//
// Required fields:
//   - role_name: the IAM role name to inspect
//   - expected_policy: the policy name expected to be attached
//
// Optional fields:
//   - endpoint: override the IAM API base URL (default: https://iam.amazonaws.com)
func checkIAMRolePolicy(fields map[string]string) (bool, string, error) {
	roleName := fields["role_name"]
	if roleName == "" {
		return false, "", fmt.Errorf("missing required field: role_name")
	}

	expectedPolicy := fields["expected_policy"]
	if expectedPolicy == "" {
		return false, "", fmt.Errorf("missing required field: expected_policy")
	}

	endpoint := fields["endpoint"]
	if endpoint == "" {
		endpoint = "https://iam.amazonaws.com"
	}

	url := fmt.Sprintf("%s/?Action=ListAttachedRolePolicies&RoleName=%s&Version=2010-05-08", endpoint, roleName)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false, "", fmt.Errorf("IAM request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("failed to read IAM response: %w", err)
	}

	var result struct {
		Policies []struct {
			PolicyName string `json:"PolicyName"`
		} `json:"AttachedPolicies"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("failed to parse IAM response: %w", err)
	}

	for _, p := range result.Policies {
		if strings.EqualFold(p.PolicyName, expectedPolicy) {
			return false, fmt.Sprintf("role %q has policy %q attached", roleName, expectedPolicy), nil
		}
	}

	return true, fmt.Sprintf("role %q is missing expected policy %q", roleName, expectedPolicy), nil
}
