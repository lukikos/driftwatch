package checker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkGlueJobStatus checks the status of an AWS Glue job via the Glue API.
// Required fields: endpoint, job_name, expected (e.g. "READY")
func checkGlueJobStatus(check config.Check) (bool, string, error) {
	endpoint, ok := check.Params["endpoint"].(string)
	if !ok || endpoint == "" {
		return false, "", fmt.Errorf("glue_job_status: missing or empty 'endpoint'")
	}

	jobName, ok := check.Params["job_name"].(string)
	if !ok || jobName == "" {
		return false, "", fmt.Errorf("glue_job_status: missing or empty 'job_name'")
	}

	expected, ok := check.Params["expected"].(string)
	if !ok || expected == "" {
		expected = "READY"
	}

	url := strings.TrimRight(endpoint, "/") + "/jobs/" + jobName
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return false, "", fmt.Errorf("glue_job_status: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return true, fmt.Sprintf("job %q not found (expected %s)", jobName, expected), nil
	}

	var result struct {
		Job struct {
			Name            string `json:"Name"`
			AllocatedCapacity int    `json:"AllocatedCapacity"`
		} `json:"Job"`
		JobStatus string `json:"JobStatus"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, "", fmt.Errorf("glue_job_status: failed to decode response: %w", err)
	}

	actual := result.JobStatus
	if actual == "" {
		actual = "READY" // default when job exists but no explicit status
	}

	if !strings.EqualFold(actual, expected) {
		return true, fmt.Sprintf("glue job %q status is %q, expected %q", jobName, actual, expected), nil
	}
	return false, "", nil
}
