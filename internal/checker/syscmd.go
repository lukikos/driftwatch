package checker

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkSysCommand runs a shell command and compares its trimmed stdout
// against the expected value defined in the check config.
//
// Config fields:
//   - command (string, required): the shell command to execute
//   - expected (string, required): expected trimmed output
func checkSysCommand(check config.Check) (bool, string, error) {
	cmd, ok := check.Config["command"]
	if !ok || strings.TrimSpace(cmd) == "" {
		return false, "", fmt.Errorf("syscmd check %q: missing or empty 'command' field", check.Name)
	}

	expected, ok := check.Config["expected"]
	if !ok {
		return false, "", fmt.Errorf("syscmd check %q: missing 'expected' field", check.Name)
	}

	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false, "", fmt.Errorf("syscmd check %q: command failed: %w", check.Name, err)
	}

	actual := strings.TrimSpace(string(out))
	if actual != strings.TrimSpace(expected) {
		msg := fmt.Sprintf("syscmd %q: expected %q, got %q", check.Name, expected, actual)
		return true, msg, nil
	}

	return false, "", nil
}
