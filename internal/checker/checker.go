package checker

import (
	"fmt"

	"github.com/driftwatch/driftwatch/internal/config"
)

// Checker dispatches config.Check entries to the appropriate check function.
type Checker struct{}

// New returns a new Checker.
func New() *Checker {
	return &Checker{}
}

// Check runs the appropriate check for the given config.Check.
// It returns (drifted bool, detail string, err error).
func (c *Checker) Check(check config.Check) (bool, string, error) {
	switch check.Type {
	case "env_var":
		return checkEnvVar(check)
	case "file_hash":
		return checkFileHash(check)
	case "file_content":
		return checkFileContent(check)
	case "file_exists":
		return checkFileExists(check)
	case "http_status":
		return checkHTTPStatus(check)
	case "port_open":
		return checkPortOpen(check)
	case "process_running":
		return checkProcessRunning(check)
	case "docker_container":
		return checkDockerContainer(check)
	case "dns_resolve":
		return checkDNSResolve(check)
	case "ssl_expiry":
		return checkSSLExpiry(check)
	case "sys_command":
		return checkSysCommand(check)
	case "dir_size":
		return checkDirSize(check)
	case "last_cron_run":
		return checkLastCronRun(check)
	case "env_file":
		return checkEnvFile(check)
	case "k8s_pod":
		return checkK8sPod(check)
	default:
		return false, "", fmt.Errorf("unknown check type: %q", check.Type)
	}
}
