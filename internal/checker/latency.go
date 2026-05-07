package checker

import (
	"fmt"
	"net/http"
	"time"
)

// checkHTTPLatency performs an HTTP GET against the given URL and reports
// drift if the response latency exceeds the configured threshold.
//
// Required fields:
//   - url:         the endpoint to probe
//   - max_ms:      maximum acceptable latency in milliseconds (integer)
//
// Optional fields:
//   - method:      HTTP method to use (default: GET)
func checkHTTPLatency(fields map[string]string) (bool, string, error) {
	url, ok := fields["url"]
	if !ok || url == "" {
		return false, "", fmt.Errorf("latency check requires 'url' field")
	}

	maxMSStr, ok := fields["max_ms"]
	if !ok || maxMSStr == "" {
		return false, "", fmt.Errorf("latency check requires 'max_ms' field")
	}

	var maxMS int
	if _, err := fmt.Sscanf(maxMSStr, "%d", &maxMS); err != nil || maxMS <= 0 {
		return false, "", fmt.Errorf("latency check: 'max_ms' must be a positive integer, got %q", maxMSStr)
	}

	method := fields["method"]
	if method == "" {
		method = http.MethodGet
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return false, "", fmt.Errorf("latency check: failed to build request: %w", err)
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return false, "", fmt.Errorf("latency check: request failed: %w", err)
	}
	resp.Body.Close()
	elapsed := time.Since(start).Milliseconds()

	if elapsed > int64(maxMS) {
		return true,
			fmt.Sprintf("latency %dms exceeds threshold %dms for %s", elapsed, maxMS, url),
			nil
	}
	return false,
		fmt.Sprintf("latency %dms within threshold %dms for %s", elapsed, maxMS, url),
		nil
}
