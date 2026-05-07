package checker

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkK8sPod checks whether a Kubernetes pod matching a given name prefix
// is running in the specified namespace using kubectl.
//
// Required fields:
//   - pod_prefix: prefix of the pod name to search for
//   - namespace:  Kubernetes namespace (defaults to "default")
//   - expected:   expected status string (defaults to "Running")
func checkK8sPod(c config.Check) (bool, string, error) {
	podPrefix, ok := c.Fields["pod_prefix"]
	if !ok || strings.TrimSpace(podPrefix) == "" {
		return false, "", fmt.Errorf("k8s_pod check %q: missing or empty 'pod_prefix' field", c.Name)
	}

	namespace, ok := c.Fields["namespace"]
	if !ok || strings.TrimSpace(namespace) == "" {
		namespace = "default"
	}

	expected, ok := c.Fields["expected"]
	if !ok || strings.TrimSpace(expected) == "" {
		expected = "Running"
	}

	ctx := context.Background()
	out, err := exec.CommandContext(ctx,
		"kubectl", "get", "pods",
		"-n", namespace,
		"--no-headers",
		"-o", "custom-columns=NAME:.metadata.name,STATUS:.status.phase",
	).Output()
	if err != nil {
		return false, "", fmt.Errorf("k8s_pod check %q: kubectl error: %w", c.Name, err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		name, status := fields[0], fields[1]
		if strings.HasPrefix(name, podPrefix) {
			if status != expected {
				return true, fmt.Sprintf("pod %q status is %q, expected %q", name, status, expected), nil
			}
			return false, fmt.Sprintf("pod %q status is %q", name, status), nil
		}
	}

	return true, fmt.Sprintf("no pod with prefix %q found in namespace %q", podPrefix, namespace), nil
}
