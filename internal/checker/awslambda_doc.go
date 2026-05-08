// Package checker provides the lambda_function check type.
//
// # lambda_function
//
// Queries an AWS Lambda function's configuration endpoint and verifies
// the function's reported State matches the expected value.
//
// ## Required fields
//
//   - function_name  — Lambda function name or ARN
//   - endpoint       — Base URL of the Lambda API
//     (e.g. https://lambda.us-east-1.amazonaws.com)
//
// ## Optional fields
//
//   - expected_state — Expected function State (default: "Active")
//
// ## Example
//
//	- name: payments-lambda-state
//	  type: lambda_function
//	  fields:
//	    function_name: payments-processor
//	    endpoint: https://lambda.us-east-1.amazonaws.com
//	    expected_state: Active
package checker
