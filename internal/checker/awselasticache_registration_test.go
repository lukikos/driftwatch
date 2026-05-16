package checker

import (
	"strings"
	"testing"

	"github.com/yourusername/driftwatch/internal/config"
)

func TestElastiCacheClusterStatus_RegistrationInDispatch(t *testing.T) {
	ch := New()
	c := config.Check{
		Name:   "cache-dispatch-test",
		Type:   "elasticache_cluster",
		Params: map[string]string{}, // missing required params — expect param error, not unknown type
	}
	_, _, err := ch.Check(c)
	if err == nil {
		t.Fatal("expected error for missing params")
	}
	if strings.Contains(err.Error(), "unknown check type") {
		t.Errorf("elasticache_cluster should be registered; got: %v", err)
	}
}

func TestElastiCacheClusterStatus_UnknownTypeStillErrors(t *testing.T) {
	ch := New()
	c := config.Check{
		Name: "bad-type",
		Type: "elasticache_cluster_nonexistent",
		Params: map[string]string{
			"endpoint":   "http://localhost",
			"cluster_id": "x",
		},
	}
	_, _, err := ch.Check(c)
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
	if !strings.Contains(err.Error(), "unknown check type") {
		t.Errorf("expected 'unknown check type' error, got: %v", err)
	}
}
