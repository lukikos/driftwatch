package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/driftwatch/internal/config"
)

func elbv2Server(state string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"State": state}) //nolint:errcheck
	}))
}

func TestCheckELBv2TargetGroupHealth_MissingEndpoint(t *testing.T) {
	chk := config.Check{Type: "elbv2_target_group_health", Params: map[string]interface{}{}}
	_, _, err := checkELBv2TargetGroupHealth(chk)
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckELBv2TargetGroupHealth_MissingTargetGroupARN(t *testing.T) {
	chk := config.Check{
		Type:   "elbv2_target_group_health",
		Params: map[string]interface{}{"endpoint": "http://localhost:9999"},
	}
	_, _, err := checkELBv2TargetGroupHealth(chk)
	if err == nil {
		t.Fatal("expected error for missing target_group_arn")
	}
}

func TestCheckELBv2TargetGroupHealth_NoDrift(t *testing.T) {
	srv := elbv2Server("healthy")
	defer srv.Close()

	chk := config.Check{
		Type: "elbv2_target_group_health",
		Params: map[string]interface{}{
			"endpoint":         srv.URL,
			"target_group_arn": "arn:aws:elasticloadbalancing:us-east-1:123:targetgroup/my-tg/abc",
			"expected":         "healthy",
		},
	}
	drifted, msg, err := checkELBv2TargetGroupHealth(chk)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Fatalf("expected no drift, got: %s", msg)
	}
}

func TestCheckELBv2TargetGroupHealth_Drift(t *testing.T) {
	srv := elbv2Server("unhealthy")
	defer srv.Close()

	chk := config.Check{
		Type: "elbv2_target_group_health",
		Params: map[string]interface{}{
			"endpoint":         srv.URL,
			"target_group_arn": "arn:aws:elasticloadbalancing:us-east-1:123:targetgroup/my-tg/abc",
			"expected":         "healthy",
		},
	}
	drifted, msg, err := checkELBv2TargetGroupHealth(chk)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Fatal("expected drift")
	}
	if msg == "" {
		t.Fatal("expected non-empty drift message")
	}
}

func TestCheckELBv2TargetGroupHealth_DefaultExpectedHealthy(t *testing.T) {
	srv := elbv2Server("healthy")
	defer srv.Close()

	chk := config.Check{
		Type: "elbv2_target_group_health",
		Params: map[string]interface{}{
			"endpoint":         srv.URL,
			"target_group_arn": "arn:aws:elasticloadbalancing:us-east-1:123:targetgroup/my-tg/abc",
			// no "expected" key — should default to "healthy"
		},
	}
	drifted, _, err := checkELBv2TargetGroupHealth(chk)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Fatal("expected no drift with default expected=healthy")
	}
}

func TestCheckELBv2TargetGroupHealth_ViaChecker(t *testing.T) {
	srv := elbv2Server("draining")
	defer srv.Close()

	c := New()
	chk := config.Check{
		Name: "elb-tg-check",
		Type: "elbv2_target_group_health",
		Params: map[string]interface{}{
			"endpoint":         srv.URL,
			"target_group_arn": "arn:aws:elasticloadbalancing:us-east-1:123:targetgroup/my-tg/abc",
			"expected":         "healthy",
		},
	}
	drifted, _, err := c.Check(chk)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Fatal("expected drift for draining state")
	}
}
