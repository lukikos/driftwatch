package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/driftwatch/internal/config"
)

func openSearchServer(status string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"DomainStatus": map[string]interface{}{
				"Status":            status,
				"Processing":        false,
				"UpgradeProcessing": false,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp) //nolint:errcheck
	}))
}

func TestCheckOpenSearchDomainStatus_MissingEndpoint(t *testing.T) {
	_, _, err := checkOpenSearchDomainStatus(config.Check{
		Params: map[string]interface{}{"domain_name": "my-domain"},
	})
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckOpenSearchDomainStatus_MissingDomainName(t *testing.T) {
	_, _, err := checkOpenSearchDomainStatus(config.Check{
		Params: map[string]interface{}{"endpoint": "http://localhost"},
	})
	if err == nil {
		t.Fatal("expected error for missing domain_name")
	}
}

func TestCheckOpenSearchDomainStatus_NoDrift(t *testing.T) {
	srv := openSearchServer("Active")
	defer srv.Close()

	drift, msg, err := checkOpenSearchDomainStatus(config.Check{
		Params: map[string]interface{}{
			"endpoint":    srv.URL,
			"domain_name": "my-domain",
			"expected":    "Active",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatalf("expected no drift, got message: %s", msg)
	}
}

func TestCheckOpenSearchDomainStatus_Drift(t *testing.T) {
	srv := openSearchServer("Processing")
	defer srv.Close()

	drift, msg, err := checkOpenSearchDomainStatus(config.Check{
		Params: map[string]interface{}{
			"endpoint":    srv.URL,
			"domain_name": "my-domain",
			"expected":    "Active",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Fatal("expected drift but got none")
	}
	if msg == "" {
		t.Fatal("expected non-empty drift message")
	}
}

func TestCheckOpenSearchDomainStatus_DefaultExpectedActive(t *testing.T) {
	srv := openSearchServer("Active")
	defer srv.Close()

	// No "expected" param — should default to "Active"
	drift, _, err := checkOpenSearchDomainStatus(config.Check{
		Params: map[string]interface{}{
			"endpoint":    srv.URL,
			"domain_name": "my-domain",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift with default expected=Active")
	}
}

func TestCheckOpenSearchDomainStatus_ViaChecker(t *testing.T) {
	srv := openSearchServer("Active")
	defer srv.Close()

	ch := New()
	drift, _, err := ch.Check(config.Check{
		Type: "opensearch_domain_status",
		Params: map[string]interface{}{
			"endpoint":    srv.URL,
			"domain_name": "prod-search",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift")
	}
}
