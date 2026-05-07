// Package checker evaluates individual drift checks defined in config.
package checker

import (
	"fmt"

	"github.com/driftwatch/driftwatch/internal/config"
)

// Checker runs drift checks against a config.Check definition.
type Checker struct{}

// New creates a new Checker.
func New() *Checker {
	return &Checker{}
}

// Run executes the check described by c and returns:
//   - drifted: true if the actual state differs from expected
//   - message: human-readable description of the drift (empty when no drift)
//   - err: non-nil if the check could not be performed
func (ch *Checker) Run(c config.Check) (bool, string, error) {
	switch c.Type {
	case "env_var":
		return checkEnvVar(c)
	case "file_hash":
		return checkFileHash(c)
	case "http_status":
		return checkHTTPStatus(c)
	case "process":
		return checkProcessRunning(c)
	case "port":
		return checkPortOpen(c)
	case "docker":
		return checkDockerContainer(c)
	case "sys_command":
		return checkSysCommand(c)
	case "dns":
		return checkDNSResolve(c)
	case "ssl_expiry":
		return checkSSLExpiry(c)
	case "file_content":
		return checkFileContent(c)
	case "file_exists":
		return checkFileExists(c)
	case "dir_size":
		return checkDirSize(c)
	case "last_cron_run":
		return checkLastCronRun(c)
	case "envfile":
		return checkEnvFile(c)
	default:
		return false, "", fmt.Errorf("unknown check type: %q", c.Type)
	}
}
