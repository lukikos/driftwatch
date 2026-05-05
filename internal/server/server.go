// Package server wires up the HTTP API for driftwatch.
package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/driftwatch/driftwatch/internal/history"
)

// Server holds the HTTP server and its dependencies.
type Server struct {
	http *http.Server
}

// New creates a new Server bound to addr, registering all API routes
// against the provided history Store.
func New(addr string, store *history.Store) *Server {
	mux := http.NewServeMux()

	mux.Handle("/status", history.StatusHandler(store))
	mux.Handle("/metrics", history.MetricsHandler(store))
	mux.Handle("/alerts", history.AlertHandler(store))

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{http: httpServer}
}

// Start begins listening and serving HTTP requests. It blocks until the
// server is shut down or encounters a fatal error.
func (s *Server) Start() error {
	log.Printf("driftwatch API listening on %s", s.http.Addr)
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the server using the provided context.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
