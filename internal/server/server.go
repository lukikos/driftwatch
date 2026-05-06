package server

import (
	"context"
	"net/http"
	"time"

	"github.com/yourusername/driftwatch/internal/history"
)

const shutdownTimeout = 5 * time.Second

// Server wraps an HTTP server and its registered routes.
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
}

// New creates a Server with all API routes registered.
func New(addr string, store *history.Store, startedAt time.Time) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", history.StatusHandler(store))
	mux.HandleFunc("/metrics", history.MetricsHandler(store))
	mux.HandleFunc("/alerts", history.AlertHandler(store))
	mux.HandleFunc("/health", history.HealthHandler(store, startedAt))

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		mux: mux,
	}
}

// Start begins listening and serving HTTP requests.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
