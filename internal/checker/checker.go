// Package checker evaluates individual infrastructure checks and reports drift.
package checker

import (
	"fmt"

	"github.com/user/driftwatch/internal/config"
)

// Checker runs configured checks and reports whether drift was detected.
type Checker struct{}

// New returns a new Checker instance.
func New() *Checker {
	return &Checker{}
}

// Check evaluates a single check and returns (drifted, message, error).
func (c *Checker) Check(check config.Check) (bool, string, error) {
	switch check.Type {
	case "env":
		return checkEnvVar(check)
	case "file_hash":
		return checkFileHash(check)
	case "http":
		return checkHTTPStatus(check)
	case "process":
		return checkProcessRunning(check)
	case "port":
		return checkPortOpen(check)
	case "docker":
		return checkDockerContainer(check)
	case "sys_command":
		return checkSysCommand(check)
	case "dns":
		return checkDNSResolve(check)
	default:
		return false, "", fmt.Errorf("unknown check type: %q", check.Type)
	}
}
