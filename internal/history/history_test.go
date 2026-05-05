package history_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/history"
)

func TestRecord_And_All(t *testing.T) {
	s := history.New(10)

	s.Record("check-env", "env_var", "VALUE changed")
	s.Record("check-file", "file_hash", "hash mismatch")

	events := s.All()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].CheckName != "check-env" {
		t.Errorf("unexpected first event name: %s", events[0].CheckName)
	}
	if events[1].CheckType != "file_hash" {
		t.Errorf("unexpected second event type: %s", events[1].CheckType)
	}
}

func TestStore_BoundedCapacity(t *testing.T) {
	s := history.New(3)

	for i := 0; i < 5; i++ {
		s.Record("c", "env_var", "drift")
	}

	if s.Len() != 3 {
		t.Errorf("expected len 3, got %d", s.Len())
	}
}

func TestStore_DefaultMaxLen(t *testing.T) {
	s := history.New(0) // should default to 100

	for i := 0; i < 150; i++ {
		s.Record("c", "env_var", "drift")
	}

	if s.Len() != 100 {
		t.Errorf("expected len 100, got %d", s.Len())
	}
}

func TestStore_Clear(t *testing.T) {
	s := history.New(10)
	s.Record("c", "env_var", "drift")
	s.Clear()

	if s.Len() != 0 {
		t.Errorf("expected empty store after Clear, got %d", s.Len())
	}
}

func TestEvent_HasTimestamp(t *testing.T) {
	s := history.New(10)
	s.Record("c", "file_hash", "mismatch")

	events := s.All()
	if events[0].DetectedAt.IsZero() {
		t.Error("expected non-zero DetectedAt timestamp")
	}
}
