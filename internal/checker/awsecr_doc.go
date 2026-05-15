// Package checker provides drift-detection logic for various infrastructure
// resources. This file documents the ecr_repository_policy check type.
//
// # ecr_repository_policy
//
// Verifies that an Amazon ECR repository's imageTagMutability setting matches
// the expected value. Drift is reported when the repository is missing or when
// the mutability setting differs from the configured expectation.
//
// ## Parameters
//
//   - endpoint    (string, required) – Base URL of the ECR-compatible API.
//     When running against real AWS, point this at the regional ECR endpoint.
//     For local testing, a LocalStack URL works well.
//   - repository  (string, required) – Name of the ECR repository to inspect.
//   - expected    (string, required) – Expected imageTagMutability value.
//     Accepted values: "MUTABLE" or "IMMUTABLE" (case-insensitive).
//
// ## Example configuration
//
//	- name: prod-ecr-immutable
//	  type: ecr_repository_policy
//	  params:
//	    endpoint: https://ecr.us-east-1.amazonaws.com
//	    repository: my-production-app
//	    expected: IMMUTABLE
package checker
