package history

import (
	"encoding/json"
	"net/http"
)

// AlertSummary represents aggregated alert statistics from the event store.
type AlertSummary struct {
	TotalEvents  int            `json:"total_events"`
	DriftCount   int            `json:"drift_count"`
	ByCheck      map[string]int `json:"by_check"`
	RecentEvents []Event        `json:"recent_events"`
}

// AlertHandler returns an HTTP handler that exposes a summary of drift alerts.
func AlertHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		events := store.All()
		byCheck := make(map[string]int)
		driftCount := 0

		for _, e := range events {
			if e.Drifted {
				driftCount++
				byCheck[e.CheckName]++
			}
		}

		recent := events
		if len(recent) > 5 {
			recent = recent[len(recent)-5:]
		}

		summary := AlertSummary{
			TotalEvents:  len(events),
			DriftCount:   driftCount,
			ByCheck:      byCheck,
			RecentEvents: recent,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(summary)
	}
}
