package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func eksServer(status string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"cluster": map[string]string{
				"status": status,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp) //nolint:errcheck
	}))
}

func TestCheckEKSClusterStatus_MissingEndpoint(t *testing.T) {
	_, _, err := checkEKSClusterStatus(config.Check{
		Fields: map[string]string{"cluster_name": "my-cluster"},
	})
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckEKSClusterStatus_MissingClusterName(t *testing.T) {
	_, _, err := checkEKSClusterStatus(config.Check{
		Fields: map[string]string{"endpoint": "http://localhost"},
	})
	if err == nil {
		t.Fatal("expected error for missing cluster_name")
	}
}

func TestCheckEKSClusterStatus_NoDrift(t *testing.T) {
	srv := eksServer("ACTIVE")
	defer srv.Close()

	drift, _, err := checkEKSClusterStatus(config.Check{
		Fields: map[string]string{
			"endpoint":     srv.URL,
			"cluster_name": "prod-cluster",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift")
	}
}

func TestCheckEKSClusterStatus_Drift(t *testing.T) {
	srv := eksServer("DEGRADED")
	defer srv.Close()

	drift, msg, err := checkEKSClusterStatus(config.Check{
		Fields: map[string]string{
			"endpoint":     srv.URL,
			"cluster_name": "prod-cluster",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Fatal("expected drift")
	}
	if msg == "" {
		t.Fatal("expected non-empty drift message")
	}
}

func TestCheckEKSClusterStatus_DefaultExpectedActive(t *testing.T) {
	srv := eksServer("ACTIVE")
	defer srv.Close()

	drift, _, err := checkEKSClusterStatus(config.Check{
		Fields: map[string]string{
			"endpoint":     srv.URL,
			"cluster_name": "prod-cluster",
			// no 'expected' field — should default to ACTIVE
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift with default expected=ACTIVE")
	}
}

func TestCheckEKSClusterStatus_CustomExpected(t *testing.T) {
	srv := eksServer("UPDATING")
	defer srv.Close()

	drift, _, err := checkEKSClusterStatus(config.Check{
		Fields: map[string]string{
			"endpoint":     srv.URL,
			"cluster_name": "prod-cluster",
			"expected":     "UPDATING",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift when status matches custom expected")
	}
}

func TestCheckEKSClusterStatus_ViaChecker(t *testing.T) {
	srv := eksServer("ACTIVE")
	defer srv.Close()

	c := New()
	drift, _, err := c.Check(config.Check{
		Name: "eks-test",
		Type: "eks_cluster_status",
		Fields: map[string]string{
			"endpoint":     srv.URL,
			"cluster_name": "prod-cluster",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift")
	}
}
