package history

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAlertHandler_EmptyStore(t *testing.T) {
	store := New(10)
	handler := AlertHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var summary AlertSummary
	if err := json.NewDecoder(rec.Body).Decode(&summary); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if summary.TotalEvents != 0 {
		t.Errorf("expected 0 total events, got %d", summary.TotalEvents)
	}
	if summary.DriftCount != 0 {
		t.Errorf("expected 0 drift count, got %d", summary.DriftCount)
	}
}

func TestAlertHandler_WithDriftEvents(t *testing.T) {
	store := New(10)
	store.Record(Event{CheckName: "file-check", Drifted: true, Timestamp: time.Now()})
	store.Record(Event{CheckName: "file-check", Drifted: true, Timestamp: time.Now()})
	store.Record(Event{CheckName: "env-check", Drifted: false, Timestamp: time.Now()})

	handler := AlertHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var summary AlertSummary
	if err := json.NewDecoder(rec.Body).Decode(&summary); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if summary.TotalEvents != 3 {
		t.Errorf("expected 3 total events, got %d", summary.TotalEvents)
	}
	if summary.DriftCount != 2 {
		t.Errorf("expected 2 drift events, got %d", summary.DriftCount)
	}
	if summary.ByCheck["file-check"] != 2 {
		t.Errorf("expected 2 for file-check, got %d", summary.ByCheck["file-check"])
	}
	if len(summary.RecentEvents) > 5 {
		t.Errorf("expected at most 5 recent events, got %d", len(summary.RecentEvents))
	}
}

func TestAlertHandler_MethodNotAllowed(t *testing.T) {
	store := New(10)
	handler := AlertHandler(store)

	req := httptest.NewRequest(http.MethodPost, "/alerts", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
