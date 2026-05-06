package checker

import (
	"fmt"
	"os/exec"
	"strings"
)

// checkDockerContainer checks whether a Docker container is running.
// Config fields:
//   - "container": name or ID of the container (required)
//   - "expected_status": expected status string, defaults to "running"
func checkDockerContainer(fields map[string]string) (bool, string, error) {
	container, ok := fields["container"]
	if !ok || container == "" {
		return false, "", fmt.Errorf("docker check missing required field: container")
	}

	expected := fields["expected_status"]
	if expected == "" {
		expected = "running"
	}

	out, err := exec.Command("docker", "inspect", "--format", "{{.State.Status}}", container).Output()
	if err != nil {
		return true, fmt.Sprintf("container %q not found or docker unavailable: %v", container, err), nil
	}

	actual := strings.TrimSpace(string(out))
	if actual != expected {
		return true, fmt.Sprintf("container %q status is %q, expected %q", container, actual, expected), nil
	}

	return false, "", nil
}
