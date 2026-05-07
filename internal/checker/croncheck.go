package checker

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkLastCronRun checks whether a cron job last ran within an expected window.
// It inspects syslog or journald output for the given job name and verifies
// the most recent execution is within max_age_minutes.
func checkLastCronRun(check config.Check) (bool, string, error) {
	jobName, ok := check.Params["job_name"]
	if !ok || strings.TrimSpace(jobName) == "" {
		return false, "", fmt.Errorf("cron_check: missing required param 'job_name'")
	}

	maxAgeStr, ok := check.Params["max_age_minutes"]
	if !ok || strings.TrimSpace(maxAgeStr) == "" {
		return false, "", fmt.Errorf("cron_check: missing required param 'max_age_minutes'")
	}

	var maxAgeMinutes int
	if _, err := fmt.Sscanf(maxAgeStr, "%d", &maxAgeMinutes); err != nil || maxAgeMinutes <= 0 {
		return false, "", fmt.Errorf("cron_check: invalid max_age_minutes %q: must be a positive integer", maxAgeStr)
	}

	// Try journalctl first, fall back to grep /var/log/syslog
	output, err := runCronLogLookup(jobName)
	if err != nil {
		return false, "", fmt.Errorf("cron_check: failed to query logs for job %q: %w", jobName, err)
	}

	if strings.TrimSpace(output) == "" {
		return true, fmt.Sprintf("no log entries found for cron job %q", jobName), nil
	}

	// Parse the last line timestamp (journalctl --since format: "Mon DD HH:MM:SS")
	lines := strings.Split(strings.TrimSpace(output), "\n")
	lastLine := lines[len(lines)-1]

	// Extract timestamp prefix (first 15 chars for syslog format)
	if len(lastLine) < 15 {
		return true, fmt.Sprintf("could not parse timestamp from log line: %q", lastLine), nil
	}

	year := time.Now().Year()
	timestamp, err := time.ParseInLocation("Jan _2 15:04:05 2006", lastLine[:15]+fmt.Sprintf(" %d", year), time.Local)
	if err != nil {
		return true, fmt.Sprintf("could not parse cron log timestamp: %v", err), nil
	}

	age := time.Since(timestamp)
	maxAge := time.Duration(maxAgeMinutes) * time.Minute

	if age > maxAge {
		return true, fmt.Sprintf("cron job %q last ran %v ago, exceeds max age of %v", jobName, age.Round(time.Second), maxAge), nil
	}

	return false, fmt.Sprintf("cron job %q last ran %v ago (within %v window)", jobName, age.Round(time.Second), maxAge), nil
}

func runCronLogLookup(jobName string) (string, error) {
	// Prefer journalctl if available
	if path, err := exec.LookPath("journalctl"); err == nil {
		cmd := exec.Command(path, "-t", "CRON", "--no-pager", "-n", "50", "--output=short")
		out, err := cmd.Output()
		if err == nil {
			return filterLines(string(out), jobName), nil
		}
	}

	// Fall back to syslog
	cmd := exec.Command("grep", "-i", jobName, "/var/log/syslog")
	out, _ := cmd.Output()
	return string(out), nil
}

func filterLines(output, needle string) string {
	var matched []string
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, needle) {
			matched = append(matched, line)
		}
	}
	return strings.Join(matched, "\n")
}
