// Package server provides the HTTP server for driftwatch's status and
// observability endpoints.
//
// It wires together all history handlers — status, metrics, alerts, health,
// and snapshot — under a single mux and exposes Start/Shutdown lifecycle
// methods for clean integration with the main daemon loop.
//
// Endpoints:
//
//	 GET /status   — full event history
//	 GET /metrics  — aggregated drift counts
//	 GET /alerts   — drift-only events
//	 GET /health   — overall health summary
//	 GET /snapshot — latest state per check (deduplicated)
package server
