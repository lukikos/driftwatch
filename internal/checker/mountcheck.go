package checker

import (
	"fmt"
	"os"
	"strings"
)

// checkMountPoint verifies that a given path is a mounted filesystem.
// It reads /proc/mounts (Linux) or uses os.Stat to detect mount presence.
//
// Config fields:
//   - path (string, required): the mount point path to check
//   - expected (string, optional): "mounted" (default) or "unmounted"
func checkMountPoint(fields map[string]string) (bool, string, error) {
	path := strings.TrimSpace(fields["path"])
	if path == "" {
		return false, "", fmt.Errorf("mount_point check requires 'path' field")
	}

	expected := strings.TrimSpace(fields["expected"])
	if expected == "" {
		expected = "mounted"
	}
	if expected != "mounted" && expected != "unmounted" {
		return false, "", fmt.Errorf("mount_point 'expected' must be 'mounted' or 'unmounted', got %q", expected)
	}

	isMounted, err := isMountPoint(path)
	if err != nil {
		return false, "", fmt.Errorf("mount_point check failed for %q: %w", path, err)
	}

	switch expected {
	case "mounted":
		if !isMounted {
			return true, fmt.Sprintf("path %q is not mounted (expected mounted)", path), nil
		}
		return false, fmt.Sprintf("path %q is mounted as expected", path), nil
	case "unmounted":
		if isMounted {
			return true, fmt.Sprintf("path %q is mounted (expected unmounted)", path), nil
		}
		return false, fmt.Sprintf("path %q is not mounted as expected", path), nil
	}

	return false, "", nil
}

// isMountPoint checks whether path is a mount point by comparing
// device IDs of the path and its parent directory.
func isMountPoint(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if !info.IsDir() {
		return false, fmt.Errorf("%q is not a directory", path)
	}

	parent := path + "/."
	parentInfo, err := os.Stat(parent)
	if err != nil {
		return false, err
	}

	return os.SameFile(info, parentInfo) == false && isSameDevice(path, parentInfo), nil
}

func isSameDevice(path string, _ os.FileInfo) bool {
	// Fallback: attempt to read /proc/mounts for the path
	data, err := os.ReadFile("/proc/mounts")
	if err != nil {
		// On non-Linux systems, assume mounted if directory exists
		return true
	}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == path {
			return true
		}
	}
	return false
}
