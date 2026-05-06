// Package checker evaluates individual infrastructure checks and reports
// whether the current state matches the expected (baseline) state.
//
// Supported check types:
//
//	file_hash – computes the SHA-256 digest of a file and compares it
//	            against the expected hash supplied in the check config.
//
//	env_var   – reads an environment variable and compares its value
//	            against the expected value supplied in the check config.
//
// New returns a Checker that can run all checks defined in the provided
// config slice. Each call to Run iterates the checks and returns a slice
// of drift results, one per check that has drifted from its baseline.
package checker
