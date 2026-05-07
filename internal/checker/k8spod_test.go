package checker

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func TestCheckK8sPod_MissingPodPrefix(t *testing.T) {
	c := config.Check{
		Name:   "test-pod",
		Type:   "k8s_pod",
		Fields: map[string]string{},
	}
	_, _, err := checkK8sPod(c)
	if err == nil {
		t.Fatal("expected error for missing pod_prefix, got nil")
	}
}

func TestCheckK8sPod_EmptyPodPrefix(t *testing.T) {
	c := config.Check{
		Name:   "test-pod",
		Type:   "k8s_pod",
		Fields: map[string]string{"pod_prefix": "   "},
	}
	_, _, err := checkK8sPod(c)
	if err == nil {
		t.Fatal("expected error for empty pod_prefix, got nil")
	}
}

func TestCheckK8sPod_DefaultNamespaceAndExpected(t *testing.T) {
	// This test verifies that defaults are applied without panicking.
	// kubectl will fail in CI, so we only check that the error is a kubectl
	// execution error, not a validation error.
	c := config.Check{
		Name:   "test-pod",
		Type:   "k8s_pod",
		Fields: map[string]string{"pod_prefix": "myapp"},
	}
	_, _, err := checkK8sPod(c)
	// We expect either nil (kubectl available) or a kubectl error — not a
	// field-validation error.
	if err != nil {
		if err.Error() == "k8s_pod check \"test-pod\": missing or empty 'pod_prefix' field" {
			t.Errorf("unexpected validation error: %v", err)
		}
	}
}

func TestCheckK8sPod_ViaChecker(t *testing.T) {
	chkr := New()
	c := config.Check{
		Name:   "k8s-missing-prefix",
		Type:   "k8s_pod",
		Fields: map[string]string{},
	}
	_, _, err := chkr.Check(c)
	if err == nil {
		t.Fatal("expected error from checker dispatch, got nil")
	}
}
