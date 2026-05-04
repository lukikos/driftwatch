// Package webhook provides a client for sending drift alert notifications
// to a configured HTTP webhook endpoint.
//
// Usage:
//
//	client := webhook.New("https://hooks.example.com/alert")
//	err := client.Send(webhook.DriftAlert{
//		CheckName: "APP_ENV",
//		CheckType: "env_var",
//		Expected:  "production",
//		Actual:    "staging",
//	})
//	if err != nil {
//		log.Printf("failed to send alert: %v", err)
//	}
//
// The DriftAlert payload is serialized as JSON and sent via HTTP POST.
// A UTC RFC3339 timestamp is automatically added to each alert.
package webhook
