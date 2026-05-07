package checker

import (
	"fmt"
	"os"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkFileExists verifies whether a file or directory exists (or does not exist)
// at the given path. Drift is reported when the actual existence state differs
// from the expected state.
//
// Required fields:
//   - path:   filesystem path to check
//
// Optional fields:
//   - expected: "true" (default) or "false" — whether the path should exist
func checkFileExists(check config.Check) (bool, string, error) {
	path, ok := check.Fields["path"]
	if !ok || path == "" {
		return false, "", fmt.Errorf("file_exists check %q: missing required field 'path'", check.Name)
	}

	expected := "true"
	if v, ok := check.Fields["expected"]; ok && v != "" {
		expected = v
	}

	if expected != "true" && expected != "false" {
		return false, "", fmt.Errorf("file_exists check %q: 'expected' must be 'true' or 'false', got %q", check.Name, expected)
	}

	_, err := os.Stat(path)
	exists := !os.IsNotExist(err)

	// Non-permission or other unexpected errors should surface.
	if err != nil && !os.IsNotExist(err) {
		return false, "", fmt.Errorf("file_exists check %q: stat error: %w", check.Name, err)
	}

	actual := "false"
	if exists {
		actual = "true"
	}

	if actual != expected {
		return true, fmt.Sprintf("path %q: expected exists=%s, got exists=%s", path, expected, actual), nil
	}

	return false, "", nil
}
