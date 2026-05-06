package history

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthResponse represents the response from the health endpoint.
type HealthResponse struct {
	Status    string    `json:"status"`
	Uptime    string    `json:"uptime"`
	Checks    int       `json:"total_checks"`
	Drifts    int       `json:"total_drifts"`
	StartedAt time.Time `json:"started_at"`
}

// HealthHandler returns an HTTP handler that reports daemon health.
// It summarises uptime and drift counts derived from the event store.
func HealthHandler(store *Store, startedAt time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		events := store.All()
		totalDrifts := 0
		for _, e := range events {
			if e.Drifted {
				totalDrifts++
			}
		}

		resp := HealthResponse{
			Status:    "ok",
			Uptime:    time.Since(startedAt).Round(time.Second).String(),
			Checks:    len(events),
			Drifts:    totalDrifts,
			StartedAt: startedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
