// Package checker provides infrastructure drift detection checks.
//
// # security_group_rules
//
// The security_group_rules check queries an AWS-compatible
// DescribeSecurityGroups endpoint and verifies that a given substring
// (such as a CIDR block or rule description) is present in the response.
//
// This is useful for detecting unintended changes to firewall rules, such
// as overly permissive ingress rules being added or expected rules being
// removed.
//
// # Configuration Fields
//
//   - endpoint (required): Base URL of the AWS-compatible security groups API.
//   - group_id (required): The security group ID to inspect (e.g. "sg-0abc123").
//   - expected (required): Substring expected to appear in the JSON response
//     (e.g. a CIDR like "10.0.0.0/8" or a description keyword).
//
// # Example
//
//	checks:
//	  - name: prod-sg-ingress
//	    type: security_group_rules
//	    fields:
//	      endpoint: "https://ec2.us-east-1.amazonaws.com"
//	      group_id: "sg-0abc1234def56789"
//	      expected: "10.0.0.0/8"
package checker
