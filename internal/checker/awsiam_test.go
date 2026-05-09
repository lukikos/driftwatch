package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func iamServer(policies []string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type policy struct {
			PolicyName string `json:"PolicyName"`
		}
		attached := make([]policy, 0, len(policies))
		for _, p := range policies {
			attached = append(attached, policy{PolicyName: p})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"AttachedPolicies": attached,
		})
	}))
}

func TestCheckIAMRolePolicy_MissingRoleName(t *testing.T) {
	_, _, err := checkIAMRolePolicy(map[string]string{
		"expected_policy": "ReadOnlyAccess",
	})
	if err == nil || err.Error() != "missing required field: role_name" {
		t.Fatalf("expected role_name error, got %v", err)
	}
}

func TestCheckIAMRolePolicy_MissingExpectedPolicy(t *testing.T) {
	_, _, err := checkIAMRolePolicy(map[string]string{
		"role_name": "my-role",
	})
	if err == nil || err.Error() != "missing required field: expected_policy" {
		t.Fatalf("expected expected_policy error, got %v", err)
	}
}

func TestCheckIAMRolePolicy_NoDrift(t *testing.T) {
	srv := iamServer([]string{"ReadOnlyAccess", "AmazonS3ReadOnlyAccess"})
	defer srv.Close()

	drift, msg, err := checkIAMRolePolicy(map[string]string{
		"role_name":       "my-role",
		"expected_policy": "ReadOnlyAccess",
		"endpoint":        srv.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatalf("expected no drift, got msg: %s", msg)
	}
}

func TestCheckIAMRolePolicy_Drift(t *testing.T) {
	srv := iamServer([]string{"AmazonS3ReadOnlyAccess"})
	defer srv.Close()

	drift, msg, err := checkIAMRolePolicy(map[string]string{
		"role_name":       "my-role",
		"expected_policy": "ReadOnlyAccess",
		"endpoint":        srv.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Fatalf("expected drift, got msg: %s", msg)
	}
}

func TestCheckIAMRolePolicy_ViaChecker(t *testing.T) {
	srv := iamServer([]string{"AdministratorAccess"})
	defer srv.Close()

	c := New()
	drift, msg, err := c.Check("iam_role_policy", map[string]string{
		"role_name":       "admin-role",
		"expected_policy": "AdministratorAccess",
		"endpoint":        srv.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatalf("expected no drift, got: %s", msg)
	}
}
