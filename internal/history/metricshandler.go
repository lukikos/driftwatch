package history

import (
	"encoding/json"
	"net/http"
)

// MetricsSummary holds aggregated drift metrics derived from the event store.
type MetricsSummary struct {
	TotalEvents  int            `json:"total_events"`
	DriftCount   int            `json:"drift_count"`
	CheckCounts  map[string]int `json:"check_counts"`
}

// MetricsHandler returns an http.HandlerFunc that serves a JSON summary
// of recorded drift events from the provided Store.
func MetricsHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		events := store.All()

		summary := MetricsSummary{
			TotalEvents: len(events),
			CheckCounts: make(map[string]int),
		}

		for _, e := range events {
			if e.Drifted {
				summary.DriftCount++
			}
			summary.CheckCounts[e.CheckName]++
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(summary); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
