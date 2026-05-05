package history_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourorg/driftwatch/internal/history"
)

func TestStatusHandler_EmptyStore(t *testing.T) {
	s := history.New(10)
	h := history.StatusHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body struct {
		Count  int              `json:"count"`
		Events []history.Event  `json:"events"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if body.Count != 0 {
		t.Errorf("expected count 0, got %d", body.Count)
	}
}

func TestStatusHandler_WithEvents(t *testing.T) {
	s := history.New(10)
	s.Record("env-check", "env_var", "PORT changed")
	s.Record("cfg-hash", "file_hash", "hash mismatch")

	h := history.StatusHandler(s)
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body struct {
		Count  int              `json:"count"`
		Events []history.Event  `json:"events"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if body.Count != 2 {
		t.Errorf("expected count 2, got %d", body.Count)
	}
	if body.Events[0].CheckName != "env-check" {
		t.Errorf("unexpected first event: %+v", body.Events[0])
	}
}

func TestStatusHandler_MethodNotAllowed(t *testing.T) {
	s := history.New(10)
	h := history.StatusHandler(s)

	req := httptest.NewRequest(http.MethodPost, "/status", nil)
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}
