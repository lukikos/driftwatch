package checker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkSymlink verifies that a symlink exists at a given path and optionally
// that it resolves to the expected target.
//
// Required fields:
//   - path: the symlink path to inspect
//
// Optional fields:
//   - expected_target: the resolved absolute path the symlink should point to
func checkSymlink(check config.Check) (bool, string, error) {
	path, ok := check.Fields["path"]
	if !ok || path == "" {
		return false, "", fmt.Errorf("symlink check requires 'path' field")
	}

	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, fmt.Sprintf("symlink does not exist at path: %s", path), nil
		}
		return false, "", fmt.Errorf("lstat %s: %w", path, err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return true, fmt.Sprintf("path exists but is not a symlink: %s", path), nil
	}

	expectedTarget, hasTarget := check.Fields["expected_target"]
	if !hasTarget || expectedTarget == "" {
		return false, fmt.Sprintf("symlink exists at %s", path), nil
	}

	resolvedTarget, err := filepath.EvalSymlinks(path)
	if err != nil {
		return false, "", fmt.Errorf("failed to resolve symlink %s: %w", path, err)
	}

	if resolvedTarget != expectedTarget {
		return true, fmt.Sprintf(
			"symlink %s resolves to %q, expected %q",
			path, resolvedTarget, expectedTarget,
		), nil
	}

	return false, fmt.Sprintf("symlink %s correctly resolves to %s", path, resolvedTarget), nil
}
