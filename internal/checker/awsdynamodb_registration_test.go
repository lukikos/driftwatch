package checker

import (
	"testing"
)

func TestDynamoDBTable_RegistrationInDispatch(t *testing.T) {
	c := New()

	// Calling with a missing required field should return an error,
	// not an "unknown check type" error — confirming the type is registered.
	_, _, err := c.Check("dynamodb_table", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing fields, got nil")
	}

	expectedSub := "missing required field"
	if !containsString(err.Error(), expectedSub) {
		t.Errorf("expected error to contain %q, got: %v", expectedSub, err)
	}
}

func TestDynamoDBTable_UnknownTypeStillErrors(t *testing.T) {
	c := New()

	_, _, err := c.Check("dynamodb_table_nonexistent", map[string]string{})
	if err == nil {
		t.Fatal("expected error for unknown check type")
	}
}

// containsString is a local helper to avoid importing strings in test files
// that already exist in this package.
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
