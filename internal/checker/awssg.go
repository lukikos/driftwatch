package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// checkSecurityGroupRules checks whether an AWS Security Group (via a mock/proxy
// endpoint) has an expected inbound rule description or cidr present.
//
// Required fields:
//   - endpoint: base URL of the AWS-compatible describe-security-groups API
//   - group_id: the security group ID to inspect
//   - expected: substring expected to appear in the JSON response body
func checkSecurityGroupRules(fields map[string]string) (bool, string, error) {
	endpoint, ok := fields["endpoint"]
	if !ok || endpoint == "" {
		return false, "", fmt.Errorf("security_group_rules: missing required field 'endpoint'")
	}

	groupID, ok := fields["group_id"]
	if !ok || groupID == "" {
		return false, "", fmt.Errorf("security_group_rules: missing required field 'group_id'")
	}

	expected, ok := fields["expected"]
	if !ok || expected == "" {
		return false, "", fmt.Errorf("security_group_rules: missing required field 'expected'")
	}

	url := fmt.Sprintf("%s?Action=DescribeSecurityGroups&GroupId=%s", endpoint, groupID)
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return false, "", fmt.Errorf("security_group_rules: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("security_group_rules: failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("security_group_rules: invalid JSON response: %w", err)
	}

	actual := string(body)
	if !contains(actual, expected) {
		return true, fmt.Sprintf("expected rule %q not found in security group %s", expected, groupID), nil
	}

	return false, "", nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		})())
}
