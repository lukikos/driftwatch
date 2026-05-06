// Package checker provides drift-detection logic for various infrastructure
// check types supported by driftwatch.
//
// Each check function accepts a map of string fields (sourced from the YAML
// configuration) and returns a (drifted bool, message string, err error) tuple.
//
// Supported check types:
//
//	env_var         – compares an environment variable to an expected value
//	file_hash       – compares a file's SHA-256 hash to an expected digest
//	file_content    – asserts a file contains a substring or matches a regex
//	http_status     – verifies an HTTP endpoint returns an expected status code
//	process_running – checks whether a named process is currently running
//	port_open       – tests whether a TCP port is accepting connections
//	docker_container– inspects a Docker container's running status
//	sys_command     – runs an arbitrary command and compares stdout to expected
//	dns_resolve     – resolves a hostname and checks for an expected IP
//	ssl_expiry      – verifies a TLS certificate has sufficient days remaining
package checker
