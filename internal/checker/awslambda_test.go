package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func lambdaServer(t *testing.T, state string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"State": state}) //nolint:errcheck
	}))
}

func TestCheckLambdaFunction_MissingFunctionName(t *testing.T) {
	_, _, err := checkLambdaFunction(config.Check{
		Fields: map[string]string{"endpoint": "http://localhost"},
	})
	if err == nil {
		t.Fatal("expected error for missing function_name")
	}
}

func TestCheckLambdaFunction_MissingEndpoint(t *testing.T) {
	_, _, err := checkLambdaFunction(config.Check{
		Fields: map[string]string{"function_name": "my-fn"},
	})
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckLambdaFunction_NoDrift(t *testing.T) {
	srv := lambdaServer(t, "Active")
	defer srv.Close()

	drift, msg, err := checkLambdaFunction(config.Check{
		Fields: map[string]string{
			"function_name": "my-fn",
			"endpoint":      srv.URL,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift, got message: %s", msg)
	}
}

func TestCheckLambdaFunction_Drift(t *testing.T) {
	srv := lambdaServer(t, "Failed")
	defer srv.Close()

	drift, msg, err := checkLambdaFunction(config.Check{
		Fields: map[string]string{
			"function_name": "my-fn",
			"endpoint":      srv.URL,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Errorf("expected drift, got message: %s", msg)
	}
}

func TestCheckLambdaFunction_DefaultExpectedActive(t *testing.T) {
	srv := lambdaServer(t, "Active")
	defer srv.Close()

	drift, _, err := checkLambdaFunction(config.Check{
		Fields: map[string]string{
			"function_name": "my-fn",
			"endpoint":      srv.URL,
			// no expected_state — should default to "Active"
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift when state is Active and default is Active")
	}
}

func TestCheckLambdaFunction_CustomExpectedState(t *testing.T) {
	srv := lambdaServer(t, "Inactive")
	defer srv.Close()

	drift, _, err := checkLambdaFunction(config.Check{
		Fields: map[string]string{
			"function_name":  "my-fn",
			"endpoint":       srv.URL,
			"expected_state": "Inactive",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift when state matches custom expected")
	}
}

func TestCheckLambdaFunction_ViaChecker(t *testing.T) {
	srv := lambdaServer(t, "Active")
	defer srv.Close()

	c := New()
	drift, _, err := c.Check(config.Check{
		Type: "lambda_function",
		Fields: map[string]string{
			"function_name": "my-fn",
			"endpoint":      srv.URL,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift")
	}
}
