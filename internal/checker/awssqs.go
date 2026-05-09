package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// checkSQSQueueAttributes checks an AWS SQS queue's attributes via the
// SQS HTTP API (or a compatible mock endpoint) and compares a specific
// attribute value against an expected value.
//
// Required fields:
//   - endpoint:  base URL of the SQS-compatible endpoint
//   - queue_url: full queue URL (used as the QueueUrl parameter)
//   - attribute: attribute name to inspect (e.g. "ApproximateNumberOfMessages")
//   - expected:  expected string value of the attribute
func checkSQSQueueAttributes(fields map[string]string) (bool, string, error) {
	endpoint := strings.TrimRight(fields["endpoint"], "/")
	if endpoint == "" {
		return false, "", fmt.Errorf("sqs_queue_attributes: missing required field 'endpoint'")
	}

	queueURL := fields["queue_url"]
	if queueURL == "" {
		return false, "", fmt.Errorf("sqs_queue_attributes: missing required field 'queue_url'")
	}

	attribute := fields["attribute"]
	if attribute == "" {
		return false, "", fmt.Errorf("sqs_queue_attributes: missing required field 'attribute'")
	}

	expected := fields["expected"]
	if expected == "" {
		return false, "", fmt.Errorf("sqs_queue_attributes: missing required field 'expected'")
	}

	url := fmt.Sprintf("%s/?Action=GetQueueAttributes&QueueUrl=%s&AttributeName.1=%s",
		endpoint, queueURL, attribute)

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return false, "", fmt.Errorf("sqs_queue_attributes: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("sqs_queue_attributes: failed to read response: %w", err)
	}

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("sqs_queue_attributes: failed to parse response: %w", err)
	}

	actual, ok := result[attribute]
	if !ok {
		return false, "", fmt.Errorf("sqs_queue_attributes: attribute %q not found in response", attribute)
	}

	if actual != expected {
		return true, fmt.Sprintf("attribute %q is %q, expected %q", attribute, actual, expected), nil
	}
	return false, "", nil
}
