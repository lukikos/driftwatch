package history_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/history"
)

// decodeHealthResponse is a helper that decodes the response body into a HealthResponse.
func decodeHealthResponse(t *testing.T, rec *httptest.ResponseRecorder) history.HealthResponse {
	t.Helper()
	var resp history.HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return resp
}

func TestHealthHandler_EmptyStore(t *testing.T) {
	store := history.New(10)
	startedAt := time.Now().Add(-5 * time.Minute)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	history.HealthHandler(store, startedAt).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	resp := decodeHealthResponse(t, rec)

	if resp.Status != "ok" {
		t.Errorf("expected status ok, got %q", resp.Status)
	}
	if resp.Checks != 0 {
		t.Errorf("expected 0 checks, got %d", resp.Checks)
	}
	if resp.Drifts != 0 {
		t.Errorf("expected 0 drifts, got %d", resp.Drifts)
	}
	if resp.Uptime == "" {
		t.Error("expected non-empty uptime")
	}
}

func TestHealthHandler_WithDriftEvents(t *testing.T) {
	store := history.New(10)
	startedAt := time.Now().Add(-2 * time.Minute)

	store.Record(history.Event{CheckName: "file-check", Drifted: true})
	store.Record(history.Event{CheckName: "env-check", Drifted: false})
	store.Record(history.Event{CheckName: "file-check", Drifted: true})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	history.HealthHandler(store, startedAt).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	resp := decodeHealthResponse(t, rec)

	if resp.Checks != 3 {
		t.Errorf("expected 3 checks, got %d", resp.Checks)
	}
	if resp.Drifts != 2 {
		t.Errorf("expected 2 drifts, got %d", resp.Drifts)
	}
}

func TestHealthHandler_MethodNotAllowed(t *testing.T) {
	store := history.New(10)
	startedAt := time.Now()

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()

	history.HealthHandler(store, startedAt).ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
