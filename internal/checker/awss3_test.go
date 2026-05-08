package checker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckS3BucketAccess_MissingURL(t *testing.T) {
	_, _, err := checkS3BucketAccess(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing url, got nil")
	}
}

func TestCheckS3BucketAccess_EmptyURL(t *testing.T) {
	_, _, err := checkS3BucketAccess(map[string]string{"url": "   "})
	if err == nil {
		t.Fatal("expected error for empty url, got nil")
	}
}

func TestCheckS3BucketAccess_NoDrift(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	drifted, msg, err := checkS3BucketAccess(map[string]string{"url": ts.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Errorf("expected no drift, got message: %s", msg)
	}
}

func TestCheckS3BucketAccess_Drift(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	drifted, msg, err := checkS3BucketAccess(map[string]string{
		"url":             ts.URL,
		"expected_status": "200",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Error("expected drift, got none")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckS3BucketAccess_DefaultExpected200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// No expected_status field — should default to "200"
	drifted, _, err := checkS3BucketAccess(map[string]string{"url": ts.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Error("expected no drift with default 200")
	}
}

func TestCheckS3BucketAccess_InvalidURL(t *testing.T) {
	drifted, msg, err := checkS3BucketAccess(map[string]string{
		"url": "http://127.0.0.1:0/no-server",
	})
	if err != nil {
		t.Fatalf("unexpected hard error: %v", err)
	}
	if !drifted {
		t.Error("expected drift for unreachable URL")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckS3BucketAccess_ViaChecker(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := New()
	drifted, _, err := c.Check("s3_bucket_access", map[string]string{
		"url": fmt.Sprintf("%s/bucket", ts.URL),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Error("expected no drift via checker dispatch")
	}
}
