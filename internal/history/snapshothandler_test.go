package history

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSnapshotHandler_EmptyStore(t *testing.T) {
	store := New(10)
	handler := SnapshotHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/snapshot", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp SnapshotResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.TotalChecks != 0 || resp.DriftCount != 0 || resp.HealthyCount != 0 {
		t.Errorf("expected empty snapshot, got %+v", resp)
	}
}

func TestSnapshotHandler_DeduplicatesByCheckName(t *testing.T) {
	store := New(10)
	now := time.Now().UTC()

	store.Record(Event{CheckName: "file-check", Drifted: false, Message: "ok", Timestamp: now.Add(-2 * time.Minute)})
	store.Record(Event{CheckName: "file-check", Drifted: true, Message: "drifted", Timestamp: now.Add(-1 * time.Minute)})
	store.Record(Event{CheckName: "env-check", Drifted: false, Message: "ok", Timestamp: now})

	handler := SnapshotHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/snapshot", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	var resp SnapshotResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if resp.TotalChecks != 2 {
		t.Errorf("expected 2 unique checks, got %d", resp.TotalChecks)
	}
	if resp.DriftCount != 1 {
		t.Errorf("expected 1 drifted, got %d", resp.DriftCount)
	}
	if resp.HealthyCount != 1 {
		t.Errorf("expected 1 healthy, got %d", resp.HealthyCount)
	}

	for _, c := range resp.Checks {
		if c.Name == "file-check" && !c.Drifted {
			t.Error("expected file-check to show latest drifted=true")
		}
	}
}

func TestSnapshotHandler_MethodNotAllowed(t *testing.T) {
	store := New(10)
	handler := SnapshotHandler(store)

	req := httptest.NewRequest(http.MethodPost, "/snapshot", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
