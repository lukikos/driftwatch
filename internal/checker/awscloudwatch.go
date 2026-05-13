package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// checkCloudWatchAlarm checks whether a CloudWatch alarm is in the expected state.
// Fields used from Check.Fields:
//   - endpoint (string, optional): override AWS endpoint URL (useful for testing)
//   - alarm_name (string, required): name of the CloudWatch alarm
//   - expected_state (string, optional): expected alarm state; defaults to "OK"
func checkCloudWatchAlarm(c Check) (bool, string, error) {
	alarmName, ok := c.Fields["alarm_name"]
	if !ok || alarmName == "" {
		return false, "", fmt.Errorf("cloudwatch_alarm: missing required field 'alarm_name'")
	}

	expectedState := "OK"
	if v, ok := c.Fields["expected_state"]; ok && v != "" {
		expectedState = v
	}

	endpoint := "https://monitoring.us-east-1.amazonaws.com"
	if v, ok := c.Fields["endpoint"]; ok && v != "" {
		endpoint = v
	}

	params := url.Values{}
	params.Set("Action", "DescribeAlarms")
	params.Set("AlarmNames.member.1", alarmName)
	params.Set("Version", "2010-08-01")

	reqURL := fmt.Sprintf("%s/?%s", endpoint, params.Encode())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return false, "", fmt.Errorf("cloudwatch_alarm: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("cloudwatch_alarm: failed to read response: %w", err)
	}

	var result struct {
		State      string `json:"state"`
		StatusCode int    `json:"status_code"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("cloudwatch_alarm: failed to parse response: %w", err)
	}

	actualState := result.State
	if result.StatusCode != 0 {
		actualState = strconv.Itoa(result.StatusCode)
	}

	if actualState != expectedState {
		return true, fmt.Sprintf("alarm %q state is %q, expected %q", alarmName, actualState, expectedState), nil
	}
	return false, fmt.Sprintf("alarm %q is in expected state %q", alarmName, expectedState), nil
}
