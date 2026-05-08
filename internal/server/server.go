package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/driftwatch/driftwatch/internal/history"
)

// Server wraps an HTTP server with registered driftwatch routes.
type Server struct {
	httpServer *http.Server
}

// New creates a new Server with all routes registered against the given store.
func New(addr string, store *history.Store) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", history.StatusHandler(store))
	mux.HandleFunc("/metrics", history.MetricsHandler(store))
	mux.HandleFunc("/alerts", history.AlertHandler(store))
	mux.HandleFunc("/health", history.HealthHandler(store))
	mux.HandleFunc("/snapshot", history.SnapshotHandler(store))

	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// Start begins listening and serving HTTP requests.
func (s *Server) Start() error {
	log.Printf("server: listening on %s", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the server with a timeout.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
