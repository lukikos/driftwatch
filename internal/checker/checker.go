// Package checker evaluates individual drift checks defined in config.
package checker

import (
	"fmt"
)

// Result holds the outcome of a single drift check.
type Result struct {
	Drifted bool
	Message string
}

// Checker dispatches named check types to their implementations.
type Checker struct{}

// New returns a new Checker.
func New() *Checker {
	return &Checker{}
}

// Check runs the named check type with the given fields.
// It returns a Result and any execution error.
func (c *Checker) Check(checkType string, fields map[string]string) (Result, error) {
	var (
		drifted bool
		msg     string
		err     error
	)

	switch checkType {
	case "env_var":
		drifted, msg, err = checkEnvVar(fields)
	case "file_hash":
		drifted, msg, err = checkFileHash(fields)
	case "http_status":
		drifted, msg, err = checkHTTPStatus(fields)
	case "process_running":
		drifted, msg, err = checkProcessRunning(fields)
	case "port_open":
		drifted, msg, err = checkPortOpen(fields)
	case "docker_container":
		drifted, msg, err = checkDockerContainer(fields)
	case "sys_command":
		drifted, msg, err = checkSysCommand(fields)
	case "dns_resolve":
		drifted, msg, err = checkDNSResolve(fields)
	case "ssl_expiry":
		drifted, msg, err = checkSSLExpiry(fields)
	case "file_content":
		drifted, msg, err = checkFileContent(fields)
	case "file_exists":
		drifted, msg, err = checkFileExists(fields)
	case "dirsize":
		drifted, msg, err = checkDirSize(fields)
	default:
		return Result{}, fmt.Errorf("unknown check type: %q", checkType)
	}

	if err != nil {
		return Result{}, err
	}

	return Result{Drifted: drifted, Message: msg}, nil
}
