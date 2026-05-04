package webhook_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/driftwatch/internal/webhook"
)

func TestSend_Success(t *testing.T) {
	var received webhook.DriftAlert

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := webhook.New(server.URL)
	alert := webhook.DriftAlert{
		CheckName: "my-check",
		CheckType: "env_var",
		Expected:  "production",
		Actual:    "staging",
	}

	if err := client.Send(alert); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if received.CheckName != "my-check" {
		t.Errorf("expected check_name 'my-check', got '%s'", received.CheckName)
	}
	if received.Timestamp == "" {
		t.Error("expected timestamp to be set")
	}
}

func TestSend_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := webhook.New(server.URL)
	err := client.Send(webhook.DriftAlert{CheckName: "test"})
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestSend_InvalidURL(t *testing.T) {
	client := webhook.New("http://127.0.0.1:0/no-server")
	err := client.Send(webhook.DriftAlert{CheckName: "test"})
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
