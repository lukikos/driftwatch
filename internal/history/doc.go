// Package history provides an in-memory store for drift detection events,
// along with HTTP handlers for querying event history, metrics, and alerts.
//
// The Store type holds a bounded ring-buffer of Event values. Each Event
// records the check name, whether drift was detected, the observed value,
// and the UTC timestamp at which the check ran.
//
// HTTP handlers expose the store over JSON:
//
//	/status  – all recorded events
//	/metrics – aggregate counts (total, drifted, clean)
//	/alerts  – only events where drift was detected
package history
