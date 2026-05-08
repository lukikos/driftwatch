package checker

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func TestCertPin_RegistrationInDispatch(t *testing.T) {
	c := New()

	// A deliberately bad host so we get a connection error (not an "unknown type" error),
	// which proves the check type is registered in the dispatch table.
	check := config.Check{
		Name: "dispatch-certpin",
		Type: "certpin",
		Params: map[string]interface{}{
			"host": "127.0.0.1",
			"port": "1", // port 1 is almost certainly closed
			"pin":  "aabbcc",
		},
	}

	_, _, err := c.Run(check)
	if err == nil {
		// connection succeeded on port 1 in some exotic CI environment — that's fine,
		// the important thing is it didn't error with "unknown check type".
		return
	}

	if err.Error() == "unknown check type: certpin" {
		t.Errorf("certpin check type is not registered in the dispatcher")
	}
}

func TestCertPin_UnknownTypeStillErrors(t *testing.T) {
	c := New()

	check := config.Check{
		Name:   "bad-type",
		Type:   "certpin_nonexistent",
		Params: map[string]interface{}{},
	}

	_, _, err := c.Run(check)
	if err == nil {
		t.Error("expected error for unknown check type")
	}
}
