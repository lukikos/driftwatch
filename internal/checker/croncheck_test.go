package checker

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func TestCheckLastCronRun_MissingJobName(t *testing.T) {
	check := config.Check{
		Name:   "cron-test",
		Type:   "cron_check",
		Params: map[string]string{"max_age_minutes": "60"},
	}
	_, _, err := checkLastCronRun(check)
	if err == nil {
		t.Fatal("expected error for missing job_name, got nil")
	}
}

func TestCheckLastCronRun_MissingMaxAge(t *testing.T) {
	check := config.Check{
		Name:   "cron-test",
		Type:   "cron_check",
		Params: map[string]string{"job_name": "backup"},
	}
	_, _, err := checkLastCronRun(check)
	if err == nil {
		t.Fatal("expected error for missing max_age_minutes, got nil")
	}
}

func TestCheckLastCronRun_InvalidMaxAge(t *testing.T) {
	check := config.Check{
		Name:   "cron-test",
		Type:   "cron_check",
		Params: map[string]string{"job_name": "backup", "max_age_minutes": "notanumber"},
	}
	_, _, err := checkLastCronRun(check)
	if err == nil {
		t.Fatal("expected error for invalid max_age_minutes, got nil")
	}
}

func TestCheckLastCronRun_ZeroMaxAge(t *testing.T) {
	check := config.Check{
		Name:   "cron-test",
		Type:   "cron_check",
		Params: map[string]string{"job_name": "backup", "max_age_minutes": "0"},
	}
	_, _, err := checkLastCronRun(check)
	if err == nil {
		t.Fatal("expected error for zero max_age_minutes, got nil")
	}
}

func TestCheckLastCronRun_NoLogEntriesReportsDrift(t *testing.T) {
	// With a job name unlikely to appear in any log, we expect either no error
	// and drift=true (no entries found), or a graceful error from log lookup.
	check := config.Check{
		Name:   "cron-test",
		Type:   "cron_check",
		Params: map[string]string{"job_name": "driftwatch_nonexistent_job_xyz", "max_age_minutes": "60"},
	}
	drift, msg, err := checkLastCronRun(check)
	if err != nil {
		// Acceptable: log lookup tool unavailable in test environment
		t.Logf("log lookup returned error (acceptable in CI): %v", err)
		return
	}
	if !drift {
		t.Errorf("expected drift=true for unknown job, got false; msg=%q", msg)
	}
}

func TestCheckLastCronRun_ViaChecker(t *testing.T) {
	c := New()
	check := config.Check{
		Name:   "cron-via-checker",
		Type:   "cron_check",
		Params: map[string]string{"job_name": "driftwatch_nonexistent_job_xyz", "max_age_minutes": "30"},
	}
	_, _, err := c.Check(check)
	if err != nil {
		t.Logf("checker returned error (acceptable in CI): %v", err)
	}
}

func TestFilterLines(t *testing.T) {
	input := "Jan  1 00:00:01 host CRON[123]: backup started\nJan  1 00:00:02 host CRON[124]: unrelated\nJan  1 00:00:03 host CRON[125]: backup done"
	result := filterLines(input, "backup")
	if result == "" {
		t.Fatal("expected filtered lines, got empty string")
	}
	for _, line := range []string{"backup started", "backup done"} {
		if !containsStr(result, line) {
			t.Errorf("expected result to contain %q", line)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
