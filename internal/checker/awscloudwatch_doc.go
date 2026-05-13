// Package checker provides infrastructure check implementations for driftwatch.
//
// # CloudWatch Alarm Check
//
// The cloudwatch_alarm check type verifies that a named AWS CloudWatch alarm
// is in the expected state (default: "OK").
//
// # Configuration Fields
//
//   - alarm_name (string, required): The name of the CloudWatch alarm to inspect.
//   - expected_state (string, optional): The expected alarm state. Defaults to "OK".
//     Common values: "OK", "ALARM", "INSUFFICIENT_DATA".
//   - endpoint (string, optional): Override the AWS endpoint URL. Useful for
//     local testing or mock servers.
//
// # Example
//
//	checks:
//	  - name: prod-cpu-alarm
//	    type: cloudwatch_alarm
//	    fields:
//	      alarm_name: high-cpu-utilization
//	      expected_state: OK
package checker
