package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ivov/driftwatch/internal/config"
)

func kinesisServer(status string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"StreamDescriptionSummary": map[string]string{
				"StreamStatus": status,
			},
		})
	}))
}

func TestCheckKinesisStreamStatus_MissingEndpoint(t *testing.T) {
	check := config.Check{
		Type:   "kinesis_stream_status",
		Params: map[string]string{"stream_name": "my-stream"},
	}
	_, _, err := checkKinesisStreamStatus(check)
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckKinesisStreamStatus_MissingStreamName(t *testing.T) {
	check := config.Check{
		Type:   "kinesis_stream_status",
		Params: map[string]string{"endpoint": "http://localhost:4566"},
	}
	_, _, err := checkKinesisStreamStatus(check)
	if err == nil {
		t.Fatal("expected error for missing stream_name")
	}
}

func TestCheckKinesisStreamStatus_NoDrift(t *testing.T) {
	srv := kinesisServer("ACTIVE")
	defer srv.Close()

	check := config.Check{
		Type: "kinesis_stream_status",
		Params: map[string]string{
			"endpoint":    srv.URL,
			"stream_name": "my-stream",
			"expected":    "ACTIVE",
		},
	}
	drifted, msg, err := checkKinesisStreamStatus(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Errorf("expected no drift, got message: %s", msg)
	}
}

func TestCheckKinesisStreamStatus_Drift(t *testing.T) {
	srv := kinesisServer("CREATING")
	defer srv.Close()

	check := config.Check{
		Type: "kinesis_stream_status",
		Params: map[string]string{
			"endpoint":    srv.URL,
			"stream_name": "my-stream",
			"expected":    "ACTIVE",
		},
	}
	drifted, msg, err := checkKinesisStreamStatus(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Error("expected drift but got none")
	}
	if msg == "" {
		t.Error("expected drift message to be non-empty")
	}
}

func TestCheckKinesisStreamStatus_DefaultExpectedActive(t *testing.T) {
	srv := kinesisServer("ACTIVE")
	defer srv.Close()

	check := config.Check{
		Type: "kinesis_stream_status",
		Params: map[string]string{
			"endpoint":    srv.URL,
			"stream_name": "my-stream",
			// no "expected" param — should default to "ACTIVE"
		},
	}
	drifted, _, err := checkKinesisStreamStatus(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Error("expected no drift with default expected=ACTIVE")
	}
}

func TestCheckKinesisStreamStatus_ViaChecker(t *testing.T) {
	srv := kinesisServer("DELETING")
	defer srv.Close()

	ch := New()
	check := config.Check{
		Name: "kinesis-test",
		Type: "kinesis_stream_status",
		Params: map[string]string{
			"endpoint":    srv.URL,
			"stream_name": "my-stream",
			"expected":    "ACTIVE",
		},
	}
	drifted, _, err := ch.Check(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Error("expected drift via checker dispatch")
	}
}
