// Package checker evaluates individual drift checks defined in configuration.
package checker

import (
	"fmt"

	"github.com/driftwatch/driftwatch/internal/config"
)

// Checker evaluates config.Check entries and reports whether drift is detected.
type Checker struct{}

// New returns a new Checker instance.
func New() *Checker {
	return &Checker{}
}

// Check runs the appropriate check function based on check.Type.
// It returns (drift bool, message string, err error).
// drift=true means the system has drifted from the expected state.
func (c *Checker) Check(check config.Check) (bool, string, error) {
	switch check.Type {
	case "env_var":
		return checkEnvVar(check)
	case "file_hash":
		return checkFileHash(check)
	case "http_status":
		return checkHTTPStatus(check)
	case "process_running":
		return checkProcessRunning(check)
	case "port_open":
		return checkPortOpen(check)
	case "docker_container":
		return checkDockerContainer(check)
	case "sys_command":
		return checkSysCommand(check)
	case "dns_resolve":
		return checkDNSResolve(check)
	case "ssl_expiry":
		return checkSSLExpiry(check)
	case "file_content":
		return checkFileContent(check)
	case "file_exists":
		return checkFileExists(check)
	case "dir_size":
		return checkDirSize(check)
	case "cron_check":
		return checkLastCronRun(check)
	default:
		return false, "", fmt.Errorf("unknown check type: %q", check.Type)
	}
}
