package checker

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// checkProcessRunning verifies that a named process is currently running on the host.
// The check config should have:
//   - name: human-readable check name
//   - process: the process name to search for (e.g. "nginx", "sshd")
func checkProcessRunning(cfg map[string]string) (bool, string, error) {
	process, ok := cfg["process"]
	if !ok || process == "" {
		return false, "", fmt.Errorf("process check requires 'process' field")
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", process))
	default:
		cmd = exec.Command("pgrep", "-x", process)
	}

	out, err := cmd.Output()
	if err != nil {
		// pgrep exits non-zero when no process found — that is drift, not an error
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return true, fmt.Sprintf("process %q is not running", process), nil
		}
		return false, "", fmt.Errorf("process check failed for %q: %w", process, err)
	}

	output := strings.TrimSpace(string(out))
	if output == "" {
		return true, fmt.Sprintf("process %q is not running", process), nil
	}

	return false, "", nil
}
