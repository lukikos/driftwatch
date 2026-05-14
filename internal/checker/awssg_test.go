package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func sgServer(t *testing.T, groupID, cidr string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qID := r.URL.Query().Get("GroupId")
		if qID != groupID {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
			"SecurityGroups": []map[string]interface{}{
				{
					"GroupId": groupID,
					"IpPermissions": []map[string]interface{}{
						{"CidrIp": cidr, "FromPort": 443},
					},
				},
			},
		})
	}))
}

func TestCheckSecurityGroupRules_MissingEndpoint(t *testing.T) {
	_, _, err := checkSecurityGroupRules(map[string]string{
		"group_id": "sg-123",
		"expected": "0.0.0.0/0",
	})
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckSecurityGroupRules_MissingGroupID(t *testing.T) {
	_, _, err := checkSecurityGroupRules(map[string]string{
		"endpoint": "http://localhost",
		"expected": "0.0.0.0/0",
	})
	if err == nil {
		t.Fatal("expected error for missing group_id")
	}
}

func TestCheckSecurityGroupRules_MissingExpected(t *testing.T) {
	_, _, err := checkSecurityGroupRules(map[string]string{
		"endpoint": "http://localhost",
		"group_id": "sg-123",
	})
	if err == nil {
		t.Fatal("expected error for missing expected")
	}
}

func TestCheckSecurityGroupRules_NoDrift(t *testing.T) {
	srv := sgServer(t, "sg-abc", "10.0.0.0/8")
	defer srv.Close()

	drift, msg, err := checkSecurityGroupRules(map[string]string{
		"endpoint": srv.URL,
		"group_id": "sg-abc",
		"expected": "10.0.0.0/8",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatalf("expected no drift, got: %s", msg)
	}
}

func TestCheckSecurityGroupRules_Drift(t *testing.T) {
	srv := sgServer(t, "sg-abc", "10.0.0.0/8")
	defer srv.Close()

	drift, msg, err := checkSecurityGroupRules(map[string]string{
		"endpoint": srv.URL,
		"group_id": "sg-abc",
		"expected": "0.0.0.0/0",
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

func TestCheckSecurityGroupRules_GroupNotFound(t *testing.T) {
	srv := sgServer(t, "sg-abc", "10.0.0.0/8")
	defer srv.Close()

	_, _, err := checkSecurityGroupRules(map[string]string{
		"endpoint": srv.URL,
		"group_id": "sg-nonexistent",
		"expected": "10.0.0.0/8",
	})
	if err == nil {
		t.Fatal("expected error for unknown group_id")
	}
}

func TestCheckSecurityGroupRules_ViaChecker(t *testing.T) {
	srv := sgServer(t, "sg-xyz", "192.168.1.0/24")
	defer srv.Close()

	c := New()
	drift, _, err := c.Check("security_group_rules", map[string]string{
		"endpoint": srv.URL,
		"group_id": "sg-xyz",
		"expected": "192.168.1.0/24",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift")
	}
}
