// Package notifier provides a rate-limited notification layer that wraps
// the webhook sender to prevent alert fatigue during repeated drift events.
package notifier

import (
	"log"
	"sync"
	"time"
)

// Sender is the interface for sending drift alerts.
type Sender interface {
	Send(checkName, message string) error
}

// Notifier wraps a Sender and suppresses duplicate alerts within a cooldown window.
type Notifier struct {
	sender   Sender
	cooldown time.Duration
	mu       sync.Mutex
	lastSent map[string]time.Time
}

// New creates a Notifier with the given sender and cooldown duration.
// If cooldown is zero, a default of 5 minutes is used.
func New(sender Sender, cooldown time.Duration) *Notifier {
	if cooldown <= 0 {
		cooldown = 5 * time.Minute
	}
	return &Notifier{
		sender:   sender,
		cooldown: cooldown,
		lastSent: make(map[string]time.Time),
	}
}

// Notify sends an alert for the given check only if the cooldown period has
// elapsed since the last alert for that check name.
func (n *Notifier) Notify(checkName, message string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if last, ok := n.lastSent[checkName]; ok {
		if time.Since(last) < n.cooldown {
			log.Printf("notifier: suppressing alert for %q (cooldown active)", checkName)
			return nil
		}
	}

	if err := n.sender.Send(checkName, message); err != nil {
		return err
	}

	n.lastSent[checkName] = time.Now()
	return nil
}

// Reset clears the cooldown state for a specific check, allowing the next
// alert to be sent immediately.
func (n *Notifier) Reset(checkName string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.lastSent, checkName)
}
