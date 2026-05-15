// Package checker provides infrastructure drift detection checks.
//
// # SSM Parameter Check
//
// The ssm_parameter check verifies that an AWS Systems Manager (SSM) Parameter
// Store parameter holds the expected value. It queries a configurable HTTP
// endpoint, making it compatible with LocalStack and other AWS-compatible
// testing environments.
//
// # Configuration Fields
//
//   - endpoint (string, required): Base URL of the SSM-compatible API.
//     Example: "http://localhost:4566"
//
//   - parameter_name (string, required): The full SSM parameter path.
//     Example: "/myapp/database/password"
//
//   - expected (string, required): The expected string value of the parameter.
//
// # Example YAML
//
//	- name: check-app-env
//	  type: ssm_parameter
//	  fields:
//	    endpoint: "http://localhost:4566"
//	    parameter_name: "/myapp/env"
//	    expected: "production"
package checker
