// Package checker provides infrastructure drift detection checks.
//
// # DynamoDB Table Check
//
// The "dynamodb_table" check verifies that an AWS DynamoDB table is in the
// expected status by querying a DynamoDB-compatible REST endpoint.
//
// ## Configuration Fields
//
//   - endpoint (required): Base URL of the DynamoDB-compatible API.
//   - table_name (required): Name of the DynamoDB table to inspect.
//   - expected (optional): Expected table status. Defaults to "ACTIVE".
//
// ## Example
//
//	checks:
//	  - name: dynamodb-users-table
//	    type: dynamodb_table
//	    fields:
//	      endpoint: http://localhost:8000
//	      table_name: users
//	      expected: ACTIVE
package checker
