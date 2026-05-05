package history

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMetricsHandler_EmptyStore(t *testing.T) {
	store := New(10)
	handler := MetricsHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var summary MetricsSummary
	if err := json.NewDecoder(rr.Body).Decode(&summary); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if summary.TotalEvents != 0 {
		t.Errorf("expected 0 total events, got %d", summary.TotalEvents)
	}
	if summary.DriftCount != 0 {
		t.Errorf("expected 0 drift count, got %d", summary.DriftCount)
	}
}

func TestMetricsHandler_WithEvents(t *testing.T) {
	store := New(10)
	store.Record(Event{CheckName: "check-a", Drifted: true, Timestamp: time.Now()})
	store.Record(Event{CheckName: "check-a", Drifted: false, Timestamp: time.Now()})
	store.Record(Event{CheckName: "check-b", Drifted: true, Timestamp: time.Now()})

	handler := MetricsHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var summary MetricsSummary
	if err := json.NewDecoder(rr.Body).Decode(&summary); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if summary.TotalEvents != 3 {
		t.Errorf("expected 3 total events, got %d", summary.TotalEvents)
	}
	if summary.DriftCount != 2 {
		t.Errorf("expected 2 drift events, got %d", summary.DriftCount)
	}
	if summary.CheckCounts["check-a"] != 2 {
		t.Errorf("expected 2 events for check-a, got %d", summary.CheckCounts["check-a"])
	}
	if summary.CheckCounts["check-b"] != 1 {
		t.Errorf("expected 1 event for check-b, got %d", summary.CheckCounts["check-b"])
	}
}

func TestMetricsHandler_MethodNotAllowed(t *testing.T) {
	store := New(10)
	handler := MetricsHandler(store)

	req := httptest.NewRequest(http.MethodPost, "/metrics", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
