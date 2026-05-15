// Package checker provides infrastructure check implementations for driftwatch.
//
// # AWS Glue Job Status Check
//
// The glue_job_status check type queries an AWS Glue-compatible endpoint to
// verify that a named Glue job is in the expected state.
//
// Required parameters:
//   - endpoint:  Base URL of the Glue API (e.g. http://localhost:4566)
//   - job_name:  Name of the Glue job to inspect
//
// Optional parameters:
//   - expected:  Expected job status string (default: "READY")
//
// Example YAML:
//
//	- name: etl-pipeline-ready
//	  type: glue_job_status
//	  params:
//	    endpoint: http://localhost:4566
//	    job_name: my-etl-job
//	    expected: READY
package checker
