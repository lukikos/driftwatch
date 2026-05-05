package history

import (
	"encoding/json"
	"net/http"
)

// StatusHandler returns an http.HandlerFunc that serves the drift event
// history as a JSON array. It is intended to be mounted on a lightweight
// diagnostic HTTP server.
func StatusHandler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		events := s.All()

		type response struct {
			Count  int     `json:"count"`
			Events []Event `json:"events"`
		}

		resp := response{
			Count:  len(events),
			Events: events,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
