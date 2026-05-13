package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func dynamoDBServer(status string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
			"Table": map[string]string{
				"TableStatus": status,
			},
		})
	}))
}

func TestCheckDynamoDBTable_MissingEndpoint(t *testing.T) {
	_, _, err := checkDynamoDBTable(map[string]string{"table_name": "users"})
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestCheckDynamoDBTable_MissingTableName(t *testing.T) {
	_, _, err := checkDynamoDBTable(map[string]string{"endpoint": "http://localhost:8000"})
	if err == nil {
		t.Fatal("expected error for missing table_name")
	}
}

func TestCheckDynamoDBTable_NoDrift(t *testing.T) {
	srv := dynamoDBServer("ACTIVE")
	defer srv.Close()

	drift, msg, err := checkDynamoDBTable(map[string]string{
		"endpoint":   srv.URL,
		"table_name": "users",
		"expected":   "ACTIVE",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift, got message: %s", msg)
	}
}

func TestCheckDynamoDBTable_Drift(t *testing.T) {
	srv := dynamoDBServer("CREATING")
	defer srv.Close()

	drift, msg, err := checkDynamoDBTable(map[string]string{
		"endpoint":   srv.URL,
		"table_name": "users",
		"expected":   "ACTIVE",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Fatal("expected drift but got none")
	}
	if msg == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestCheckDynamoDBTable_DefaultExpectedActive(t *testing.T) {
	srv := dynamoDBServer("ACTIVE")
	defer srv.Close()

	drift, _, err := checkDynamoDBTable(map[string]string{
		"endpoint":   srv.URL,
		"table_name": "orders",
		// no "expected" field — should default to ACTIVE
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift with default expected=ACTIVE")
	}
}

func TestCheckDynamoDBTable_ViaChecker(t *testing.T) {
	srv := dynamoDBServer("DELETING")
	defer srv.Close()

	c := New()
	drift, msg, err := c.Check("dynamodb_table", map[string]string{
		"endpoint":   srv.URL,
		"table_name": "archive",
		"expected":   "ACTIVE",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Fatal("expected drift via checker dispatch")
	}
	if msg == "" {
		t.Error("expected drift message via checker dispatch")
	}
}
