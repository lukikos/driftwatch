package checker

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkEnvFile checks that a .env file contains an expected key=value pair.
// Fields used:
//   - path: path to the .env file (required)
//   - key: the variable name to look up (required)
//   - expected: the expected value for the key (required)
func checkEnvFile(c config.Check) (bool, string, error) {
	path, ok := c.Fields["path"]
	if !ok || path == "" {
		return false, "", fmt.Errorf("envfile check %q: missing required field 'path'", c.Name)
	}

	key, ok := c.Fields["key"]
	if !ok || key == "" {
		return false, "", fmt.Errorf("envfile check %q: missing required field 'key'", c.Name)
	}

	expected, ok := c.Fields["expected"]
	if !ok {
		return false, "", fmt.Errorf("envfile check %q: missing required field 'expected'", c.Name)
	}

	f, err := os.Open(path)
	if err != nil {
		return false, "", fmt.Errorf("envfile check %q: cannot open file: %w", c.Name, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		if strings.TrimSpace(parts[0]) == key {
			actual := strings.TrimSpace(parts[1])
			if actual == expected {
				return false, "", nil
			}
			return true, fmt.Sprintf("key %q: expected %q, got %q", key, expected, actual), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, "", fmt.Errorf("envfile check %q: error reading file: %w", c.Name, err)
	}

	return true, fmt.Sprintf("key %q not found in %s", key, path), nil
}
