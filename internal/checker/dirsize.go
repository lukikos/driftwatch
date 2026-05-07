package checker

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// checkDirSize checks whether the total size of a directory (in bytes)
// is within an expected threshold. Fields:
//   - path: directory to measure
//   - max_bytes: maximum allowed size in bytes (as string)
func checkDirSize(fields map[string]string) (bool, string, error) {
	path, ok := fields["path"]
	if !ok || path == "" {
		return false, "", fmt.Errorf("dirsize check requires 'path' field")
	}

	maxStr, ok := fields["max_bytes"]
	if !ok || maxStr == "" {
		return false, "", fmt.Errorf("dirsize check requires 'max_bytes' field")
	}

	maxBytes, err := strconv.ParseInt(maxStr, 10, 64)
	if err != nil {
		return false, "", fmt.Errorf("dirsize: invalid max_bytes %q: %w", maxStr, err)
	}

	var total int64
	err = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			total += info.Size()
		}
		return nil
	})
	if err != nil {
		return false, "", fmt.Errorf("dirsize: failed to walk %q: %w", path, err)
	}

	if total > maxBytes {
		return true, fmt.Sprintf("directory %q size %d bytes exceeds max %d bytes", path, total, maxBytes), nil
	}

	return false, fmt.Sprintf("directory %q size %d bytes within limit %d bytes", path, total, maxBytes), nil
}
