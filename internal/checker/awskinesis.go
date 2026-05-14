package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ivov/driftwatch/internal/config"
)

// checkKinesisStreamStatus checks whether an AWS Kinesis stream has the expected status.
// It uses a configurable endpoint (for testability) and queries the DescribeStreamSummary API.
//
// Required fields:
//   - endpoint:      base URL of the Kinesis API (e.g. http://localhost:4566)
//   - stream_name:   name of the Kinesis stream to check
//   - expected:      expected stream status (default: "ACTIVE")
func checkKinesisStreamStatus(check config.Check) (bool, string, error) {
	endpoint, ok := check.Params["endpoint"]
	if !ok || strings.TrimSpace(endpoint) == "" {
		return false, "", fmt.Errorf("kinesis_stream_status: missing required param 'endpoint'")
	}

	streamName, ok := check.Params["stream_name"]
	if !ok || strings.TrimSpace(streamName) == "" {
		return false, "", fmt.Errorf("kinesis_stream_status: missing required param 'stream_name'")
	}

	expected := check.Params["expected"]
	if strings.TrimSpace(expected) == "" {
		expected = "ACTIVE"
	}

	url := fmt.Sprintf("%s/?Action=DescribeStreamSummary&StreamName=%s", strings.TrimRight(endpoint, "/"), streamName)

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return false, "", fmt.Errorf("kinesis_stream_status: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("kinesis_stream_status: failed to read response: %w", err)
	}

	var result struct {
		StreamDescriptionSummary struct {
			StreamStatus string `json:"StreamStatus"`
		} `json:"StreamDescriptionSummary"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("kinesis_stream_status: failed to parse response: %w", err)
	}

	actual := result.StreamDescriptionSummary.StreamStatus
	if actual != expected {
		return true, fmt.Sprintf("stream %q status is %q, expected %q", streamName, actual, expected), nil
	}

	return false, "", nil
}
