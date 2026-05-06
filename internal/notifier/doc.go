// Package notifier provides a rate-limited alert dispatcher for driftwatch.
//
// It wraps any Sender (typically a webhook client) and ensures that repeated
// drift detections for the same check do not flood downstream systems. Each
// check name has its own independent cooldown window; once the window expires
// the next drift event will trigger a fresh alert.
//
// # Architecture
//
// The Notifier sits between the drift-detection layer and the delivery layer:
//
//	┌─────────────┐     ┌──────────────┐     ┌────────────┐
//	│  Detector   │────▶│   Notifier   │────▶│   Sender   │
//	└─────────────┘     └──────────────┘     └────────────┘
//
// The Notifier maintains an in-memory map of last-alert timestamps keyed by
// check name. Calls to Notify that fall within the cooldown window are silently
// dropped; no error is returned to the caller.
//
// # Basic usage
//
//	wh := webhook.New(cfg.WebhookURL)
//	n  := notifier.New(wh, 5*time.Minute)
//	n.Notify(check.Name, "value changed")
//
// # Concurrency
//
// Notifier is safe for concurrent use by multiple goroutines.
package notifier
