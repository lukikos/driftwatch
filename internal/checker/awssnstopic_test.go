package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/driftwatch/internal/config"
)

func snsServer(attrs map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
			"Attributes": attrs,
		})
	}))
}

func TestCheckSNSTopicAttributes_MissingEndpoint(t *testing.T) {
	_, _, err := checkSNSTopicAttributes(config.Check{Params: map[string]string{
		"topic_arn": "arn:aws:sns:us-east-1:123456789012:MyTopic",
		"attribute": "DisplayName",
		"expected":  "MyTopic",
	}})
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckSNSTopicAttributes_MissingTopicARN(t *testing.T) {
	_, _, err := checkSNSTopicAttributes(config.Check{Params: map[string]string{
		"endpoint":  "http://localhost:4566",
		"attribute": "DisplayName",
		"expected":  "MyTopic",
	}})
	if err == nil {
		t.Fatal("expected error for missing topic_arn")
	}
}

func TestCheckSNSTopicAttributes_MissingAttribute(t *testing.T) {
	_, _, err := checkSNSTopicAttributes(config.Check{Params: map[string]string{
		"endpoint":  "http://localhost:4566",
		"topic_arn": "arn:aws:sns:us-east-1:123456789012:MyTopic",
		"expected":  "MyTopic",
	}})
	if err == nil {
		t.Fatal("expected error for missing attribute")
	}
}

func TestCheckSNSTopicAttributes_MissingExpected(t *testing.T) {
	_, _, err := checkSNSTopicAttributes(config.Check{Params: map[string]string{
		"endpoint":  "http://localhost:4566",
		"topic_arn": "arn:aws:sns:us-east-1:123456789012:MyTopic",
		"attribute": "DisplayName",
	}})
	if err == nil {
		t.Fatal("expected error for missing expected")
	}
}

func TestCheckSNSTopicAttributes_NoDrift(t *testing.T) {
	srv := snsServer(map[string]string{"DisplayName": "MyTopic", "SubscriptionsConfirmed": "3"})
	defer srv.Close()

	drift, _, err := checkSNSTopicAttributes(config.Check{Params: map[string]string{
		"endpoint":  srv.URL,
		"topic_arn": "arn:aws:sns:us-east-1:123456789012:MyTopic",
		"attribute": "DisplayName",
		"expected":  "MyTopic",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift")
	}
}

func TestCheckSNSTopicAttributes_Drift(t *testing.T) {
	srv := snsServer(map[string]string{"DisplayName": "WrongName"})
	defer srv.Close()

	drift, msg, err := checkSNSTopicAttributes(config.Check{Params: map[string]string{
		"endpoint":  srv.URL,
		"topic_arn": "arn:aws:sns:us-east-1:123456789012:MyTopic",
		"attribute": "DisplayName",
		"expected":  "MyTopic",
	}})
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

func TestCheckSNSTopicAttributes_ViaChecker(t *testing.T) {
	srv := snsServer(map[string]string{"SubscriptionsConfirmed": "5"})
	defer srv.Close()

	c := New()
	drift, _, err := c.Check(config.Check{
		Type: "sns_topic_attributes",
		Params: map[string]string{
			"endpoint":  srv.URL,
			"topic_arn": "arn:aws:sns:us-east-1:123456789012:MyTopic",
			"attribute": "SubscriptionsConfirmed",
			"expected":  "5",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Fatal("expected no drift")
	}
}
