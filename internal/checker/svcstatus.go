package checker

import (
	"fmt"
	"os/exec"
	"strings"
)

// checkServiceStatus checks whether a systemd service is in the expected state.
// Required fields:
//   - service: name of the systemd service (e.g. "nginx")
//   - expected: expected ActiveState value (e.g. "active", "inactive", "failed")
func checkServiceStatus(fields map[string]string) (bool, string, error) {
	service, ok := fields["service"]
	if !ok || strings.TrimSpace(service) == "" {
		return false, "", fmt.Errorf("svc_status check missing required field: service")
	}

	expected, ok := fields["expected"]
	if !ok || strings.TrimSpace(expected) == "" {
		expected = "active"
	}

	out, err := exec.Command("systemctl", "show", "-p", "ActiveState", "--value", service).Output()
	if err != nil {
		return false, "", fmt.Errorf("svc_status: failed to query service %q: %w", service, err)
	}

	actual := strings.TrimSpace(string(out))
	if actual != expected {
		return true, fmt.Sprintf("service %q is %q, expected %q", service, actual, expected), nil
	}

	return false, "", nil
}
