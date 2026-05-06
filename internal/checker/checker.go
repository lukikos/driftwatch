// Package checker evaluates individual drift checks defined in config.
package checker

import "fmt"

// Checker runs named drift checks.
type Checker struct{}

// New creates a new Checker.
func New() *Checker {
	return &Checker{}
}

// Check dispatches a check by type and returns (drifted, message, error).
func (c *Checker) Check(checkType string, fields map[string]string) (bool, string, error) {
	switch checkType {
	case "env_var":
		return checkEnvVar(fields)
	case "file_hash":
		return checkFileHash(fields)
	case "http_status":
		return checkHTTPStatus(fields)
	case "process_running":
		return checkProcessRunning(fields)
	case "port_open":
		return checkPortOpen(fields)
	case "docker_container":
		return checkDockerContainer(fields)
	default:
		return false, "", fmt.Errorf("unknown check type: %q", checkType)
	}
}
