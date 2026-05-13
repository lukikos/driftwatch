package checker

import (
	"strings"
	"testing"
)

func TestCloudWatchAlarm_RegistrationInDispatch(t *testing.T) {
	chkr := New()
	c := Check{
		Name: "cw-registration",
		Type: "cloudwatch_alarm",
		Fields: map[string]string{
			// missing alarm_name triggers a known error, not "unknown check type"
		},
	}
	_, _, err := chkr.Run(c)
	if err == nil {
		t.Fatal("expected error for missing alarm_name field")
	}
	if strings.Contains(err.Error(), "unknown check type") {
		t.Errorf("cloudwatch_alarm should be registered; got: %v", err)
	}
}

func TestCloudWatchAlarm_UnknownTypeStillErrors(t *testing.T) {
	chkr := New()
	c := Check{
		Name: "bad-type",
		Type: "cloudwatch_nonexistent",
		Fields: map[string]string{},
	}
	_, _, err := chkr.Run(c)
	if err == nil {
		t.Fatal("expected error for unknown check type")
	}
	if !strings.Contains(err.Error(), "unknown check type") {
		t.Errorf("expected 'unknown check type' error, got: %v", err)
	}
}
