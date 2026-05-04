package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DriftAlert represents the payload sent to the webhook endpoint.
type DriftAlert struct {
	CheckName string `json:"check_name"`
	CheckType string `json:"check_type"`
	Expected  string `json:"expected"`
	Actual    string `json:"actual"`
	Timestamp string `json:"timestamp"`
}

// Client sends drift alerts to a configured webhook URL.
type Client struct {
	URL        string
	HTTPClient *http.Client
}

// New creates a new webhook Client with the given URL.
func New(url string) *Client {
	return &Client{
		URL: url,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send serializes the DriftAlert and POSTs it to the webhook URL.
func (c *Client) Send(alert DriftAlert) error {
	alert.Timestamp = time.Now().UTC().Format(time.RFC3339)

	payload, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("webhook: failed to marshal alert: %w", err)
	}

	resp, err := c.HTTPClient.Post(c.URL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("webhook: failed to send alert: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
