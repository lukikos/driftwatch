package checker

import (
	"fmt"
	"strconv"
	"syscall"
)

// checkDiskUsage checks whether disk usage percentage on a given path
// exceeds a configured threshold. Fields:
//   - path: filesystem path to check (required)
//   - max_percent: maximum allowed usage percentage, e.g. "85" (required)
func checkDiskUsage(fields map[string]string) (bool, string, error) {
	path, ok := fields["path"]
	if !ok || path == "" {
		return false, "", fmt.Errorf("disk_usage check requires 'path' field")
	}

	maxStr, ok := fields["max_percent"]
	if !ok || maxStr == "" {
		return false, "", fmt.Errorf("disk_usage check requires 'max_percent' field")
	}

	maxPercent, err := strconv.ParseFloat(maxStr, 64)
	if err != nil {
		return false, "", fmt.Errorf("disk_usage: invalid max_percent %q: %w", maxStr, err)
	}
	if maxPercent <= 0 || maxPercent > 100 {
		return false, "", fmt.Errorf("disk_usage: max_percent must be between 1 and 100, got %v", maxPercent)
	}

	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return false, "", fmt.Errorf("disk_usage: statfs %q failed: %w", path, err)
	}

	total := stat.Blocks * uint64(stat.Bsize)
	if total == 0 {
		return false, "", fmt.Errorf("disk_usage: total size of %q is zero", path)
	}

	avail := stat.Bavail * uint64(stat.Bsize)
	used := total - avail
	usedPercent := float64(used) / float64(total) * 100.0

	msg := fmt.Sprintf("disk usage at %q is %.1f%% (max %.1f%%)", path, usedPercent, maxPercent)
	if usedPercent > maxPercent {
		return true, msg, nil
	}
	return false, msg, nil
}
