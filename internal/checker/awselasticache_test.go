package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/driftwatch/internal/config"
)

func elastiCacheServer(status string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": status})
	}))
}

func TestCheckElastiCacheClusterStatus_MissingEndpoint(t *testing.T) {
	c := config.Check{Type: "elasticache_cluster", Params: map[string]string{"cluster_id": "my-cluster"}}
	_, _, err := checkElastiCacheClusterStatus(c)
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckElastiCacheClusterStatus_MissingClusterID(t *testing.T) {
	c := config.Check{Type: "elasticache_cluster", Params: map[string]string{"endpoint": "http://localhost"}}
	_, _, err := checkElastiCacheClusterStatus(c)
	if err == nil {
		t.Fatal("expected error for missing cluster_id")
	}
}

func TestCheckElastiCacheClusterStatus_NoDrift(t *testing.T) {
	srv := elastiCacheServer("available")
	defer srv.Close()

	c := config.Check{
		Type: "elasticache_cluster",
		Params: map[string]string{
			"endpoint":   srv.URL,
			"cluster_id": "prod-cache",
			"expected":   "available",
		},
	}
	drift, msg, err := checkElastiCacheClusterStatus(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift, got message: %s", msg)
	}
}

func TestCheckElastiCacheClusterStatus_Drift(t *testing.T) {
	srv := elastiCacheServer("modifying")
	defer srv.Close()

	c := config.Check{
		Type: "elasticache_cluster",
		Params: map[string]string{
			"endpoint":   srv.URL,
			"cluster_id": "prod-cache",
			"expected":   "available",
		},
	}
	drift, msg, err := checkElastiCacheClusterStatus(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Error("expected drift but got none")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckElastiCacheClusterStatus_DefaultExpectedAvailable(t *testing.T) {
	srv := elastiCacheServer("available")
	defer srv.Close()

	c := config.Check{
		Type: "elasticache_cluster",
		Params: map[string]string{
			"endpoint":   srv.URL,
			"cluster_id": "prod-cache",
			// no expected — should default to "available"
		},
	}
	drift, _, err := checkElastiCacheClusterStatus(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift with default expected=available")
	}
}

func TestCheckElastiCacheClusterStatus_ViaChecker(t *testing.T) {
	srv := elastiCacheServer("available")
	defer srv.Close()

	ch := New()
	c := config.Check{
		Name: "cache-check",
		Type: "elasticache_cluster",
		Params: map[string]string{
			"endpoint":   srv.URL,
			"cluster_id": "prod-cache",
		},
	}
	drift, _, err := ch.Check(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift")
	}
}
