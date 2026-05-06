// Package checker provides infrastructure check implementations for driftwatch.
//
// Supported check types:
//
//   - env_var:    compares an environment variable against an expected value
//   - file_hash:  computes the SHA-256 hash of a file and compares it to an expected digest
//   - http_status: performs an HTTP GET and validates the response status code
//   - process:    verifies that a named process is running via /proc or `pgrep`
//   - port:       attempts a TCP connection to host:port to confirm it is open
//   - docker:     inspects a Docker container's status via the Docker socket
//   - syscmd:     executes an arbitrary shell command and compares trimmed stdout
//                 to an expected value; useful for ad-hoc or platform-specific checks
//
// Each check function accepts a config.Check and returns (drifted bool, message string, err error).
package checker
