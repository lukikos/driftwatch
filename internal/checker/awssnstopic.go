package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/yourusername/driftwatch/internal/config"
)

// checkSNSTopicAttributes checks an AWS SNS topic attribute against an expected value.
// It supports a custom endpoint for testing/localstack scenarios.
//
// Required fields:
//   - endpoint:       base URL of the SNS-compatible API
//   - topic_arn:      the ARN of the SNS topic
//   - attribute:      attribute name to inspect (e.g. "DisplayName", "SubscriptionsConfirmed")
//   - expected:       expected value of the attribute
func checkSNSTopicAttributes(check config.Check) (bool, string, error) {
	endpoint, ok := check.Params["endpoint"]
	if !ok || strings.TrimSpace(endpoint) == "" {
		return false, "", fmt.Errorf("sns_topic_attributes: missing or empty 'endpoint'")
	}

	topicARN, ok := check.Params["topic_arn"]
	if !ok || strings.TrimSpace(topicARN) == "" {
		return false, "", fmt.Errorf("sns_topic_attributes: missing or empty 'topic_arn'")
	}

	attribute, ok := check.Params["attribute"]
	if !ok || strings.TrimSpace(attribute) == "" {
		return false, "", fmt.Errorf("sns_topic_attributes: missing or empty 'attribute'")
	}

	expected, ok := check.Params["expected"]
	if !ok {
		return false, "", fmt.Errorf("sns_topic_attributes: missing 'expected'")
	}

	url := fmt.Sprintf("%s/?Action=GetTopicAttributes&TopicArn=%s", strings.TrimRight(endpoint, "/"), topicARN)

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return false, "", fmt.Errorf("sns_topic_attributes: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("sns_topic_attributes: failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("sns_topic_attributes: failed to parse response: %w", err)
	}

	attrs, ok := result["Attributes"].(map[string]interface{})
	if !ok {
		return false, "", fmt.Errorf("sns_topic_attributes: 'Attributes' key missing or malformed in response")
	}

	actual, exists := attrs[attribute]
	if !exists {
		return false, "", fmt.Errorf("sns_topic_attributes: attribute %q not found in response", attribute)
	}

	actualStr := fmt.Sprintf("%v", actual)
	if actualStr != expected {
		return true, fmt.Sprintf("attribute %q: expected %q, got %q", attribute, expected, actualStr), nil
	}

	return false, "", nil
}
