// Package checker evaluates individual drift checks.
package checker

import "fmt"

// Checker dispatches named check types to their implementations.
type Checker struct{}

// New returns a new Checker.
func New() *Checker {
	return &Checker{}
}

// Check runs the named check type with the provided fields.
// Returns (drifted, message, error).
func (c *Checker) Check(checkType string, fields map[string]string) (bool, string, error) {
	switch checkType {
	case "env_var":
		return checkEnvVar(fields)
	case "file_hash":
		return checkFileHash(fields)
	case "file_content":
		return checkFileContent(fields)
	case "http_status":
		return checkHTTPStatus(fields)
	case "process_running":
		return checkProcessRunning(fields)
	case "port_open":
		return checkPortOpen(fields)
	case "docker_container":
		return checkDockerContainer(fields)
	case "sys_command":
		return checkSysCommand(fields)
	case "dns_resolve":
		return checkDNSResolve(fields)
	case "ssl_expiry":
		return checkSSLExpiry(fields)
	default:
		return false, "", fmt.Errorf("unknown check type: %q", checkType)
	}
}
