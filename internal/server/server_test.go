package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/history"
)

func TestNew_RegistersRoutes(t *testing.T) {
	store := history.New(10)
	srv := New(":0", store)

	routes := []string{"/status", "/metrics", "/alerts", "/health", "/snapshot"}
	for _, route := range routes {
		req := httptest.NewRequest(http.MethodGet, route, nil)
		rr := httptest.NewRecorder()
		srv.httpServer.Handler.ServeHTTP(rr, req)
		if rr.Code == http.StatusNotFound {
			t.Errorf("route %s not registered (got 404)", route)
		}
	}
}

func TestServer_StartAndShutdown(t *testing.T) {
	store := history.New(10)
	srv := New("127.0.0.1:0", store)

	errCh := make(chan error, 1)
	go func() { errCh <- srv.Start() }()

	time.Sleep(50 * time.Millisecond)
	if err := srv.Shutdown(); err != nil {
		t.Fatalf("shutdown error: %v", err)
	}
}

func TestServer_AlertsEndpoint(t *testing.T) {
	store := history.New(10)
	store.Record(history.Event{CheckName: "cfg", Drifted: true, Message: "mismatch", Timestamp: time.Now()})

	srv := New(":0", store)
	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	rr := httptest.NewRecorder()
	srv.httpServer.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var result []map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 alert, got %d", len(result))
	}
}

func TestServer_SnapshotEndpoint(t *testing.T) {
	store := history.New(10)
	store.Record(history.Event{CheckName: "env", Drifted: false, Message: "ok", Timestamp: time.Now()})

	srv := New(":0", store)
	req := httptest.NewRequest(http.MethodGet, "/snapshot", nil)
	rr := httptest.NewRecorder()
	srv.httpServer.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp history.SnapshotResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.TotalChecks != 1 {
		t.Errorf("expected 1 check in snapshot, got %d", resp.TotalChecks)
	}
}
