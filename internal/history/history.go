// Package history provides a simple in-memory store for tracking
// drift events detected during the lifetime of the daemon.
package history

import (
	"sync"
	"time"
)

// Event represents a single detected drift occurrence.
type Event struct {
	CheckName string
	CheckType string
	Message   string
	DetectedAt time.Time
}

// Store holds a bounded ring-buffer of drift events.
type Store struct {
	mu     sync.RWMutex
	events []Event
	maxLen int
}

// New creates a Store that retains at most maxLen events.
func New(maxLen int) *Store {
	if maxLen <= 0 {
		maxLen = 100
	}
	return &Store{maxLen: maxLen}
}

// Record appends a new drift event to the store.
func (s *Store) Record(name, checkType, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e := Event{
		CheckName:  name,
		CheckType:  checkType,
		Message:    message,
		DetectedAt: time.Now().UTC(),
	}
	s.events = append(s.events, e)
	if len(s.events) > s.maxLen {
		s.events = s.events[len(s.events)-s.maxLen:]
	}
}

// All returns a snapshot of all stored events, oldest first.
func (s *Store) All() []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snap := make([]Event, len(s.events))
	copy(snap, s.events)
	return snap
}

// Len returns the number of events currently stored.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.events)
}

// Clear removes all stored events.
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = s.events[:0]
}
