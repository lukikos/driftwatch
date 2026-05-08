package checker

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// checkS3BucketAccess checks whether an S3-compatible bucket endpoint is
// accessible and returns the expected HTTP status. It performs a HEAD request
// against the bucket URL so no credentials or AWS SDK are required.
func checkS3BucketAccess(fields map[string]string) (bool, string, error) {
	url, ok := fields["url"]
	if !ok || strings.TrimSpace(url) == "" {
		return false, "", fmt.Errorf("s3_bucket_access: missing required field 'url'")
	}

	expected := fields["expected_status"]
	if expected == "" {
		expected = "200"
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Head(url)
	if err != nil {
		return true, fmt.Sprintf("request failed: %v", err), nil
	}
	defer resp.Body.Close()

	actual := fmt.Sprintf("%d", resp.StatusCode)
	if actual != expected {
		return true,
			fmt.Sprintf("bucket %s returned status %s, expected %s", url, actual, expected),
			nil
	}
	return false, "", nil
}
