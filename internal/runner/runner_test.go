package runner_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/config"
	"github.com/example/driftwatch/internal/runner"
)

func TestRun_NoDrift(t *testing.T) {
	// Webhook server should NOT be called when there is no drift.
	webhookCalled := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webhookCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	_ = os.Setenv("DW_TEST_VAR", "expected")
	defer os.Unsetenv("DW_TEST_VAR")

	cfg := &config.Config{
		WebhookURL: ts.URL,
		Interval:   50 * time.Millisecond,
		Checks: []config.Check{
			{Name: "env-ok", Type: "env_var", Key: "DW_TEST_VAR", Expected: "expected"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	runner.New(cfg).Run(ctx)

	if webhookCalled {
		t.Error("webhook should not be called when no drift detected")
	}
}

func TestRun_DriftTriggersWebhook(t *testing.T) {
	type payload struct {
		Check  string `json:"check"`
		Detail string `json:"detail"`
	}

	received := make(chan payload, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p payload
		_ = json.NewDecoder(r.Body).Decode(&p)
		received <- p
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	_ = os.Setenv("DW_TEST_VAR", "actual")
	defer os.Unsetenv("DW_TEST_VAR")

	cfg := &config.Config{
		WebhookURL: ts.URL,
		Interval:   200 * time.Millisecond,
		Checks: []config.Check{
			{Name: "env-drift", Type: "env_var", Key: "DW_TEST_VAR", Expected: "expected"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	runner.New(cfg).Run(ctx)

	select {
	case p := <-received:
		if p.Check != "env-drift" {
			t.Errorf("unexpected check name %q", p.Check)
		}
	default:
		t.Error("expected webhook to be called for drifted check")
	}
}
