package history

import (
	"encoding/json"
	"net/http"
	"time"
)

// SnapshotResponse represents the current state snapshot of all tracked checks.
type SnapshotResponse struct {
	GeneratedAt  time.Time      `json:"generated_at"`
	TotalChecks  int            `json:"total_checks"`
	DriftCount   int            `json:"drift_count"`
	HealthyCount int            `json:"healthy_count"`
	Checks       []CheckSummary `json:"checks"`
}

// CheckSummary holds the latest status for a single check.
type CheckSummary struct {
	Name     string    `json:"name"`
	Drifted  bool      `json:"drifted"`
	LastSeen time.Time `json:"last_seen"`
	Message  string    `json:"message"`
}

// SnapshotHandler returns an HTTP handler that provides a current state snapshot
// derived from the most recent event per check name in the store.
func SnapshotHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		events := store.All()
		latest := make(map[string]Event)
		for _, e := range events {
			if existing, ok := latest[e.CheckName]; !ok || e.Timestamp.After(existing.Timestamp) {
				latest[e.CheckName] = e
			}
		}

		summaries := make([]CheckSummary, 0, len(latest))
		driftCount := 0
		for _, e := range latest {
			if e.Drifted {
				driftCount++
			}
			summaries = append(summaries, CheckSummary{
				Name:     e.CheckName,
				Drifted:  e.Drifted,
				LastSeen: e.Timestamp,
				Message:  e.Message,
			})
		}

		resp := SnapshotResponse{
			GeneratedAt:  time.Now().UTC(),
			TotalChecks:  len(latest),
			DriftCount:   driftCount,
			HealthyCount: len(latest) - driftCount,
			Checks:       summaries,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
