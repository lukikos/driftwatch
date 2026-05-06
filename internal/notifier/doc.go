// Package notifier provides a rate-limited alert dispatcher for driftwatch.
//
// It wraps any Sender (typically a webhook client) and ensures that repeated
// drift detections for the same check do not flood downstream systems. Each
// check name has its own independent cooldown window; once the window expires
// the next drift event will trigger a fresh alert.
//
// Basic usage:
//
//	wh := webhook.New(cfg.WebhookURL)
//	n  := notifier.New(wh, 5*time.Minute)
//	n.Notify(check.Name, "value changed")
package notifier
