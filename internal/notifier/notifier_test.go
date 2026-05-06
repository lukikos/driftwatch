package notifier_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/notifier"
)

// mockSender records how many times Send was called.
type mockSender struct {
	callCount atomic.Int32
	failNext   bool
}

func (m *mockSender) Send(checkName, message string) error {
	if m.failNext {
		return errors.New("send failed")
	}
	m.callCount.Add(1)
	return nil
}

func TestNotify_SendsOnFirstCall(t *testing.T) {
	sender := &mockSender{}
	n := notifier.New(sender, time.Minute)

	if err := n.Notify("check-env", "drift detected"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := sender.callCount.Load(); got != 1 {
		t.Errorf("expected 1 send, got %d", got)
	}
}

func TestNotify_SuppressesDuringCooldown(t *testing.T) {
	sender := &mockSender{}
	n := notifier.New(sender, time.Hour)

	_ = n.Notify("check-env", "first alert")
	_ = n.Notify("check-env", "second alert") // should be suppressed

	if got := sender.callCount.Load(); got != 1 {
		t.Errorf("expected 1 send (cooldown suppressed second), got %d", got)
	}
}

func TestNotify_SendsAfterCooldownExpires(t *testing.T) {
	sender := &mockSender{}
	n := notifier.New(sender, 10*time.Millisecond)

	_ = n.Notify("check-file", "first")
	time.Sleep(20 * time.Millisecond)
	_ = n.Notify("check-file", "second")

	if got := sender.callCount.Load(); got != 2 {
		t.Errorf("expected 2 sends after cooldown, got %d", got)
	}
}

func TestNotify_IndependentCooldownsPerCheck(t *testing.T) {
	sender := &mockSender{}
	n := notifier.New(sender, time.Hour)

	_ = n.Notify("check-a", "drift")
	_ = n.Notify("check-b", "drift")

	if got := sender.callCount.Load(); got != 2 {
		t.Errorf("expected 2 sends for different checks, got %d", got)
	}
}

func TestNotify_ReturnsErrorFromSender(t *testing.T) {
	sender := &mockSender{failNext: true}
	n := notifier.New(sender, time.Minute)

	if err := n.Notify("check-env", "drift"); err == nil {
		t.Error("expected error from sender, got nil")
	}
}

func TestReset_AllowsImmediateResend(t *testing.T) {
	sender := &mockSender{}
	n := notifier.New(sender, time.Hour)

	_ = n.Notify("check-env", "first")
	n.Reset("check-env")
	_ = n.Notify("check-env", "second")

	if got := sender.callCount.Load(); got != 2 {
		t.Errorf("expected 2 sends after reset, got %d", got)
	}
}

func TestNew_DefaultCooldown(t *testing.T) {
	sender := &mockSender{}
	// passing zero should fall back to 5-minute default without panicking
	n := notifier.New(sender, 0)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
