package checker

import (
	"fmt"
	"net/http"
	"time"
)

// httpClient is used for HTTP checks; can be overridden in tests.
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// checkHTTPStatus performs an HTTP GET against the given URL and verifies
// the response status code matches the expected value.
func checkHTTPStatus(url string, expectedStatus int) (bool, string, error) {
	if url == "" {
		return false, "", fmt.Errorf("http_status check requires a non-empty url")
	}
	if expectedStatus == 0 {
		expectedStatus = http.StatusOK
	}

	resp, err := httpClient.Get(url) //nolint:noctx
	if err != nil {
		return false, "", fmt.Errorf("http_status check request failed: %w", err)
	}
	defer resp.Body.Close()

	actual := resp.StatusCode
	if actual != expectedStatus {
		return true, fmt.Sprintf("expected status %d, got %d", expectedStatus, actual), nil
	}
	return false, "", nil
}
