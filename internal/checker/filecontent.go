package checker

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// checkFileContent checks whether a file's content matches an expected string or regex pattern.
// Fields:
//   - path:     path to the file to read
//   - contains: substring that must be present (optional)
//   - matches:  regex pattern that must match (optional)
//
// At least one of contains or matches must be provided.
func checkFileContent(fields map[string]string) (bool, string, error) {
	path, ok := fields["path"]
	if !ok || path == "" {
		return false, "", fmt.Errorf("file_content check requires 'path' field")
	}

	contains := fields["contains"]
	pattern := fields["matches"]

	if contains == "" && pattern == "" {
		return false, "", fmt.Errorf("file_content check requires 'contains' or 'matches' field")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return false, "", fmt.Errorf("file_content: could not read file %q: %w", path, err)
	}

	content := string(data)

	if contains != "" {
		if !strings.Contains(content, contains) {
			return true, fmt.Sprintf("file %q does not contain %q", path, contains), nil
		}
	}

	if pattern != "" {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return false, "", fmt.Errorf("file_content: invalid regex %q: %w", pattern, err)
		}
		if !re.MatchString(content) {
			return true, fmt.Sprintf("file %q does not match pattern %q", path, pattern), nil
		}
	}

	return false, "", nil
}
