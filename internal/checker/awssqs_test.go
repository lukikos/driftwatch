package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func sqsServer(t *testing.T, attribute, value string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{attribute: value}) //nolint:errcheck
	}))
}

func TestCheckSQSQueueAttributes_MissingEndpoint(t *testing.T) {
	_, _, err := checkSQSQueueAttributes(map[string]string{
		"queue_url": "https://sqs.us-east-1.amazonaws.com/123/my-queue",
		"attribute": "ApproximateNumberOfMessages",
		"expected":  "0",
	})
	if err == nil || !strings.Contains(err.Error(), "endpoint") {
		t.Fatalf("expected endpoint error, got %v", err)
	}
}

func TestCheckSQSQueueAttributes_MissingQueueURL(t *testing.T) {
	_, _, err := checkSQSQueueAttributes(map[string]string{
		"endpoint":  "http://localhost",
		"attribute": "ApproximateNumberOfMessages",
		"expected":  "0",
	})
	if err == nil || !strings.Contains(err.Error(), "queue_url") {
		t.Fatalf("expected queue_url error, got %v", err)
	}
}

func TestCheckSQSQueueAttributes_MissingAttribute(t *testing.T) {
	_, _, err := checkSQSQueueAttributes(map[string]string{
		"endpoint":  "http://localhost",
		"queue_url": "https://sqs.us-east-1.amazonaws.com/123/my-queue",
		"expected":  "0",
	})
	if err == nil || !strings.Contains(err.Error(), "attribute") {
		t.Fatalf("expected attribute error, got %v", err)
	}
}

func TestCheckSQSQueueAttributes_MissingExpected(t *testing.T) {
	_, _, err := checkSQSQueueAttributes(map[string]string{
		"endpoint":  "http://localhost",
		"queue_url": "https://sqs.us-east-1.amazonaws.com/123/my-queue",
		"attribute": "ApproximateNumberOfMessages",
	})
	if err == nil || !strings.Contains(err.Error(), "expected") {
		t.Fatalf("expected 'expected' field error, got %v", err)
	}
}

func TestCheckSQSQueueAttributes_NoDrift(t *testing.T) {
	srv := sqsServer(t, "ApproximateNumberOfMessages", "0")
	defer srv.Close()

	drifted, msg, err := checkSQSQueueAttributes(map[string]string{
		"endpoint":  srv.URL,
		"queue_url": "https://sqs.us-east-1.amazonaws.com/123/my-queue",
		"attribute": "ApproximateNumberOfMessages",
		"expected":  "0",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Fatalf("expected no drift, got: %s", msg)
	}
}

func TestCheckSQSQueueAttributes_Drift(t *testing.T) {
	srv := sqsServer(t, "ApproximateNumberOfMessages", "42")
	defer srv.Close()

	drifted, msg, err := checkSQSQueueAttributes(map[string]string{
		"endpoint":  srv.URL,
		"queue_url": "https://sqs.us-east-1.amazonaws.com/123/my-queue",
		"attribute": "ApproximateNumberOfMessages",
		"expected":  "0",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Fatal("expected drift, got none")
	}
	if !strings.Contains(msg, "42") {
		t.Fatalf("expected msg to contain actual value '42', got: %s", msg)
	}
}

func TestCheckSQSQueueAttributes_ViaChecker(t *testing.T) {
	srv := sqsServer(t, "ApproximateNumberOfMessages", "0")
	defer srv.Close()

	c := New()
	drifted, _, err := c.Check("sqs_queue_attributes", map[string]string{
		"endpoint":  srv.URL,
		"queue_url": "https://sqs.us-east-1.amazonaws.com/123/my-queue",
		"attribute": "ApproximateNumberOfMessages",
		"expected":  "0",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Fatal("expected no drift")
	}
}
