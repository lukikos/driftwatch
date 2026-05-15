package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func ecrServer(t *testing.T, repo string, mutability string, statusCode int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if statusCode == http.StatusNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"imageTagMutability": mutability,
		})
	}))
}

func TestCheckECRRepositoryPolicy_MissingEndpoint(t *testing.T) {
	_, _, err := checkECRRepositoryPolicy(config.Check{
		Params: map[string]interface{}{"repository": "my-repo", "expected": "IMMUTABLE"},
	})
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckECRRepositoryPolicy_MissingRepository(t *testing.T) {
	_, _, err := checkECRRepositoryPolicy(config.Check{
		Params: map[string]interface{}{"endpoint": "http://localhost", "expected": "IMMUTABLE"},
	})
	if err == nil {
		t.Fatal("expected error for missing repository")
	}
}

func TestCheckECRRepositoryPolicy_MissingExpected(t *testing.T) {
	_, _, err := checkECRRepositoryPolicy(config.Check{
		Params: map[string]interface{}{"endpoint": "http://localhost", "repository": "my-repo"},
	})
	if err == nil {
		t.Fatal("expected error for missing expected")
	}
}

func TestCheckECRRepositoryPolicy_NoDrift(t *testing.T) {
	srv := ecrServer(t, "my-repo", "IMMUTABLE", http.StatusOK)
	defer srv.Close()

	drift, msg, err := checkECRRepositoryPolicy(config.Check{
		Params: map[string]interface{}{
			"endpoint":   srv.URL,
			"repository": "my-repo",
			"expected":   "IMMUTABLE",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift, got message: %s", msg)
	}
}

func TestCheckECRRepositoryPolicy_Drift(t *testing.T) {
	srv := ecrServer(t, "my-repo", "MUTABLE", http.StatusOK)
	defer srv.Close()

	drift, msg, err := checkECRRepositoryPolicy(config.Check{
		Params: map[string]interface{}{
			"endpoint":   srv.URL,
			"repository": "my-repo",
			"expected":   "IMMUTABLE",
		},
	})
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

func TestCheckECRRepositoryPolicy_NotFound(t *testing.T) {
	srv := ecrServer(t, "missing-repo", "", http.StatusNotFound)
	defer srv.Close()

	drift, msg, err := checkECRRepositoryPolicy(config.Check{
		Params: map[string]interface{}{
			"endpoint":   srv.URL,
			"repository": "missing-repo",
			"expected":   "IMMUTABLE",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Error("expected drift for missing repository")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckECRRepositoryPolicy_ViaChecker(t *testing.T) {
	srv := ecrServer(t, "prod-repo", "IMMUTABLE", http.StatusOK)
	defer srv.Close()

	ch := New()
	drift, _, err := ch.Check(config.Check{
		Type: "ecr_repository_policy",
		Params: map[string]interface{}{
			"endpoint":   srv.URL,
			"repository": "prod-repo",
			"expected":   "IMMUTABLE",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift")
	}
}
