package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func jsonServer(t *testing.T, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
}

func TestCheckJSONField_NoDrift(t *testing.T) {
	srv := jsonServer(t, map[string]interface{}{"status": "ok"})
	defer srv.Close()

	drift, _, err := checkJSONField(map[string]string{
		"url": srv.URL, "field": "status", "expected": "ok",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift")
	}
}

func TestCheckJSONField_Drift(t *testing.T) {
	srv := jsonServer(t, map[string]interface{}{"status": "degraded"})
	defer srv.Close()

	drift, msg, err := checkJSONField(map[string]string{
		"url": srv.URL, "field": "status", "expected": "ok",
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

func TestCheckJSONField_NestedField(t *testing.T) {
	srv := jsonServer(t, map[string]interface{}{"db": map[string]interface{}{"connected": "true"}})
	defer srv.Close()

	drift, _, err := checkJSONField(map[string]string{
		"url": srv.URL, "field": "db.connected", "expected": "true",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift")
	}
}

func TestCheckJSONField_MissingURL(t *testing.T) {
	_, _, err := checkJSONField(map[string]string{"field": "status", "expected": "ok"})
	if err == nil {
		t.Fatal("expected error for missing url")
	}
}

func TestCheckJSONField_MissingField(t *testing.T) {
	_, _, err := checkJSONField(map[string]string{"url": "http://x", "expected": "ok"})
	if err == nil {
		t.Fatal("expected error for missing field")
	}
}

func TestCheckJSONField_MissingExpected(t *testing.T) {
	_, _, err := checkJSONField(map[string]string{"url": "http://x", "field": "status"})
	if err == nil {
		t.Fatal("expected error for missing expected")
	}
}

func TestCheckJSONField_FieldNotInResponse(t *testing.T) {
	srv := jsonServer(t, map[string]interface{}{"other": "value"})
	defer srv.Close()

	_, _, err := checkJSONField(map[string]string{
		"url": srv.URL, "field": "status", "expected": "ok",
	})
	if err == nil {
		t.Fatal("expected error for missing field in response")
	}
}

func TestCheckJSONField_ViaChecker(t *testing.T) {
	srv := jsonServer(t, map[string]interface{}{"health": "healthy"})
	defer srv.Close()

	c := New()
	drift, _, err := c.Check("jsonfield", map[string]string{
		"url": srv.URL, "field": "health", "expected": "healthy",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift")
	}
}
