// Package checker provides infrastructure drift detection checks.
//
// # ELBv2 Target Group Health Check
//
// The elbv2_target_group_health check queries an AWS Application or Network
// Load Balancer target group and compares its health state against an expected
// value.
//
// # Configuration
//
//	- name:   friendly name for this check
//	  type:   elbv2_target_group_health
//	  params:
//	    endpoint:         "https://elasticloadbalancing.us-east-1.amazonaws.com"
//	    target_group_arn: "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/my-tg/abcdef"
//	    expected:         "healthy"   # optional, default: healthy
//
// # Notes
//
//   - The endpoint field allows pointing at a local mock (e.g. LocalStack)
//     for testing without live AWS credentials.
//   - Drift is reported when the actual state does not match the expected value
//     (case-insensitive comparison).
package checker
