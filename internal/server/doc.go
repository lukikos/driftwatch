// Package server provides the HTTP API server for driftwatch.
//
// It exposes three endpoints:
//
//	/status   - returns all recorded drift events as JSON
//	/metrics  - returns Prometheus-compatible plain-text metrics
//	/alerts   - returns a summary of drift alerts grouped by check name
//
// Usage:
//
//	store := history.New(100)
//	srv := server.New(":8080", store)
//	if err := srv.Start(); err != nil {
//	    log.Fatal(err)
//	}
package server
