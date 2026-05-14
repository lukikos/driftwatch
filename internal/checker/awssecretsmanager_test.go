package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/driftwatch/internal/config"
)

func secretsManagerServer(t *testing.T, secretID string, payload map[string]interface{}, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.URL.Query().Get("secretId")
		if got != secretID {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestCheckSecretsManagerSecret_MissingEndpoint(t *testing.T) {
	_, _, err := checkSecretsManagerSecret(config.Check{Params: map[string]string{"secret_id": "my-secret"}})
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckSecretsManagerSecret_MissingSecretID(t *testing.T) {
	_, _, err := checkSecretsManagerSecret(config.Check{Params: map[string]string{"endpoint": "http://localhost"}})
	if err == nil {
		t.Fatal("expected error for missing secret_id")
	}
}

func TestCheckSecretsManagerSecret_NoDrift(t *testing.T) {
	srv := secretsManagerServer(t, "prod/db", map[string]interface{}{"ARN": "arn:aws:secretsmanager:us-east-1:123:secret:prod/db"}, http.StatusOK)
	defer srv.Close()

	drift, msg, err := checkSecretsManagerSecret(config.Check{
		Params: map[string]string{"endpoint": srv.URL, "secret_id": "prod/db"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatalf("expected no drift, got message: %s", msg)
	}
}

func TestCheckSecretsManagerSecret_Drift_SecretMissing(t *testing.T) {
	srv := secretsManagerServer(t, "prod/db", nil, http.StatusNotFound)
	defer srv.Close()

	drift, msg, err := checkSecretsManagerSecret(config.Check{
		Params: map[string]string{"endpoint": srv.URL, "secret_id": "prod/db"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Fatal("expected drift for missing secret")
	}
	if msg == "" {
		t.Fatal("expected non-empty drift message")
	}
}

func TestCheckSecretsManagerSecret_NoDrift_FieldMatch(t *testing.T) {
	srv := secretsManagerServer(t, "prod/db", map[string]interface{}{"RotationEnabled": "true"}, http.StatusOK)
	defer srv.Close()

	drift, _, err := checkSecretsManagerSecret(config.Check{
		Params: map[string]string{
			"endpoint":  srv.URL,
			"secret_id": "prod/db",
			"field":     "RotationEnabled",
			"expected":  "true",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift when field matches")
	}
}

func TestCheckSecretsManagerSecret_Drift_FieldMismatch(t *testing.T) {
	srv := secretsManagerServer(t, "prod/db", map[string]interface{}{"RotationEnabled": "false"}, http.StatusOK)
	defer srv.Close()

	drift, msg, err := checkSecretsManagerSecret(config.Check{
		Params: map[string]string{
			"endpoint":  srv.URL,
			"secret_id": "prod/db",
			"field":     "RotationEnabled",
			"expected":  "true",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Fatal("expected drift for field mismatch")
	}
	if msg == "" {
		t.Fatal("expected non-empty drift message")
	}
}

func TestCheckSecretsManagerSecret_ViaChecker(t *testing.T) {
	srv := secretsManagerServer(t, "my-secret", map[string]interface{}{"Name": "my-secret"}, http.StatusOK)
	defer srv.Close()

	c := New()
	drift, _, err := c.Check(config.Check{
		Type:   "secrets_manager_secret",
		Params: map[string]string{"endpoint": srv.URL, "secret_id": "my-secret"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift")
	}
}
