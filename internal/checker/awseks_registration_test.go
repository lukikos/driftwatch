package checker

import (
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func TestEKSClusterStatus_RegistrationInDispatch(t *testing.T) {
	// Verifies that eks_cluster_status is wired into the checker dispatch table.
	// We pass an empty endpoint so the check returns an error rather than making
	// a real network call — what matters is that it is NOT treated as an unknown type.
	c := New()
	_, _, err := c.Check(config.Check{
		Name:   "eks-dispatch-test",
		Type:   "eks_cluster_status",
		Fields: map[string]string{},
	})
	if err == nil {
		t.Fatal("expected an error due to missing fields, not nil")
	}
	if strings.Contains(err.Error(), "unknown check type") {
		t.Fatalf("eks_cluster_status is not registered in the dispatch table: %v", err)
	}
}

func TestEKSClusterStatus_UnknownTypeStillErrors(t *testing.T) {
	c := New()
	_, _, err := c.Check(config.Check{
		Name: "bad-type",
		Type: "eks_cluster_status_NONEXISTENT",
	})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
	if !strings.Contains(err.Error(), "unknown check type") {
		t.Fatalf("expected 'unknown check type' error, got: %v", err)
	}
}
