// Package checker provides drift detection logic for various infrastructure
// components.
//
// # EC2 Instance Metadata Check
//
// The ec2_instance_metadata check queries the AWS EC2 Instance Metadata Service
// (IMDS) and compares a metadata key's value against an expected string.
//
// Example configuration:
//
//	checks:
//	  - name: instance-type-check
//	    type: ec2_instance_metadata
//	    fields:
//	      metadata_key: instance-type
//	      expected: t3.micro
//
// An optional metadata_url field overrides the default IMDS endpoint, which is
// useful in tests or environments using IMDSv2 token proxies.
package checker
