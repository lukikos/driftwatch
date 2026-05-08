package checker

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheckEC2InstanceMetadata_MissingKey(t *testing.T) {
	_, _, err := checkEC2InstanceMetadata(map[string]string{
		"expected": "t3.micro",
	})
	if err == nil || !strings.Contains(err.Error(), "metadata_key") {
		t.Fatalf("expected error about missing metadata_key, got %v", err)
	}
}

func TestCheckEC2InstanceMetadata_MissingExpected(t *testing.T) {
	_, _, err := checkEC2InstanceMetadata(map[string]string{
		"metadata_key": "instance-type",
	})
	if err == nil || !strings.Contains(err.Error(), "expected") {
		t.Fatalf("expected error about missing expected, got %v", err)
	}
}

func TestCheckEC2InstanceMetadata_NoDrift(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("t3.micro"))
	}))
	defer ts.Close()

	drifted, msg, err := checkEC2InstanceMetadata(map[string]string{
		"metadata_key": "instance-type",
		"expected":     "t3.micro",
		"metadata_url": ts.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Fatalf("expected no drift, got message: %s", msg)
	}
}

func TestCheckEC2InstanceMetadata_Drift(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("m5.large"))
	}))
	defer ts.Close()

	drifted, msg, err := checkEC2InstanceMetadata(map[string]string{
		"metadata_key": "instance-type",
		"expected":     "t3.micro",
		"metadata_url": ts.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted {
		t.Fatal("expected drift but got none")
	}
	if !strings.Contains(msg, "m5.large") || !strings.Contains(msg, "t3.micro") {
		t.Fatalf("unexpected drift message: %s", msg)
	}
}

func TestCheckEC2InstanceMetadata_InvalidURL(t *testing.T) {
	_, _, err := checkEC2InstanceMetadata(map[string]string{
		"metadata_key": "instance-type",
		"expected":     "t3.micro",
		"metadata_url": "http://127.0.0.1:0",
	})
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}

func TestCheckEC2InstanceMetadata_ViaChecker(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ami-0abcdef1234567890"))
	}))
	defer ts.Close()

	c := New()
	drifted, msg, err := c.Check("ec2_instance_metadata", map[string]string{
		"metadata_key": "ami-id",
		"expected":     "ami-0abcdef1234567890",
		"metadata_url": ts.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted {
		t.Fatalf("expected no drift, got: %s", msg)
	}
}
