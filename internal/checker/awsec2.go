package checker

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// checkEC2InstanceMetadata checks an EC2 instance metadata key against an expected value.
// It uses the IMDSv1 endpoint (http://169.254.169.254/latest/meta-data/) by default,
// but a custom metadata_url can be provided for testing or IMDSv2 proxies.
//
// Required fields:
//   - metadata_key: the metadata path (e.g. "instance-type")
//   - expected:     the expected value
//
// Optional fields:
//   - metadata_url: base URL (default: http://169.254.169.254/latest/meta-data)
func checkEC2InstanceMetadata(fields map[string]string) (bool, string, error) {
	key := fields["metadata_key"]
	if key == "" {
		return false, "", fmt.Errorf("ec2_instance_metadata: missing required field 'metadata_key'")
	}

	expected := fields["expected"]
	if expected == "" {
		return false, "", fmt.Errorf("ec2_instance_metadata: missing required field 'expected'")
	}

	baseURL := fields["metadata_url"]
	if baseURL == "" {
		baseURL = "http://169.254.169.254/latest/meta-data"
	}

	url := strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(key, "/")

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false, "", fmt.Errorf("ec2_instance_metadata: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("ec2_instance_metadata: failed to read response: %w", err)
	}

	actual := strings.TrimSpace(string(body))
	if actual != expected {
		return true, fmt.Sprintf("ec2 metadata %q: expected %q, got %q", key, expected, actual), nil
	}

	return false, "", nil
}
