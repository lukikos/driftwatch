package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/history"
)

func TestNew_RegistersRoutes(t *testing.T) {
	store := history.New(10)
	srv := New("127.0.0.1:0", store)

	if srv == nil {
		t.Fatal("expected non-nil server")
	}
	if srv.http == nil {
		t.Fatal("expected non-nil http.Server")
	}
}

func TestServer_StartAndShutdown(t *testing.T) {
	store := history.New(10)
	srv := New("127.0.0.1:18765", store)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	// Give the server a moment to start.
	time.Sleep(50 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown error: %v", err)
	}

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("unexpected server error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server did not stop in time")
	}
}

func TestServer_AlertsEndpoint(t *testing.T) {
	store := history.New(10)
	srv := New("127.0.0.1:18766", store)

	go srv.Start() //nolint:errcheck
	time.Sleep(50 * time.Millisecond)
	defer srv.Shutdown(context.Background()) //nolint:errcheck

	resp, err := http.Get("http://127.0.0.1:18766/alerts")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
