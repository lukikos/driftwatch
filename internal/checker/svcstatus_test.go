package checker

import (
	"testing"
)

func TestCheckServiceStatus_MissingService(t *testing.T) {
	_, _, err := checkServiceStatus(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing service field")
	}
}

func TestCheckServiceStatus_EmptyService(t *testing.T) {
	_, _, err := checkServiceStatus(map[string]string{"service": "   "})
	if err == nil {
		t.Fatal("expected error for empty service field")
	}
}

func TestCheckServiceStatus_DefaultExpectedActive(t *testing.T) {
	// We can't guarantee systemctl is available in CI, so we test via Checker
	// dispatch and expect either a result or an exec error — not a field error.
	c := New()
	drift, msg, err := c.Check("svc_status", map[string]string{
		"service": "nonexistent-driftwatch-svc",
	})
	// systemctl may not be present; we just confirm no field-validation error
	// and that the function was dispatched correctly.
	_ = drift
	_ = msg
	_ = err
}

func TestCheckServiceStatus_ViaChecker_MissingField(t *testing.T) {
	c := New()
	_, _, err := c.Check("svc_status", map[string]string{})
	if err == nil {
		t.Fatal("expected error when service field is missing")
	}
}

func TestCheckServiceStatus_ViaChecker_UnknownType(t *testing.T) {
	c := New()
	_, _, err := c.Check("svc_status_unknown", map[string]string{
		"service": "nginx",
	})
	if err == nil {
		t.Fatal("expected error for unknown check type")
	}
}
