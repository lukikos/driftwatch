// Package checker provides infrastructure drift detection for various check types.
//
// # svc_status check
//
// The svc_status check queries a systemd service's ActiveState using
// `systemctl show` and compares it to an expected value.
//
// Required fields:
//   - service: the name of the systemd service to inspect (e.g. "nginx", "sshd")
//
// Optional fields:
//   - expected: the expected ActiveState value (default: "active")
//     Common values: "active", "inactive", "failed", "activating"
//
// Example configuration:
//
//	checks:
//	  - name: nginx-running
//	    type: svc_status
//	    fields:
//	      service: nginx
//	      expected: active
package checker
