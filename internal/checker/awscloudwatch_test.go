package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func cloudWatchServer(state string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"state": state})
	}))
}

func TestCheckCloudWatchAlarm_MissingAlarmName(t *testing.T) {
	_, _, err := checkCloudWatchAlarm(Check{Fields: map[string]string{}})
	if err == nil {
		t.Fatal("expected error for missing alarm_name")
	}
}

func TestCheckCloudWatchAlarm_NoDrift(t *testing.T) {
	srv := cloudWatchServer("OK")
	defer srv.Close()

	c := Check{
		Fields: map[string]string{
			"alarm_name":     "test-alarm",
			"expected_state": "OK",
			"endpoint":       srv.URL,
		},
	}
	drift, msg, err := checkCloudWatchAlarm(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift, got: %s", msg)
	}
}

func TestCheckCloudWatchAlarm_Drift(t *testing.T) {
	srv := cloudWatchServer("ALARM")
	defer srv.Close()

	c := Check{
		Fields: map[string]string{
			"alarm_name":     "test-alarm",
			"expected_state": "OK",
			"endpoint":       srv.URL,
		},
	}
	drift, msg, err := checkCloudWatchAlarm(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Errorf("expected drift, got msg: %s", msg)
	}
}

func TestCheckCloudWatchAlarm_DefaultExpectedOK(t *testing.T) {
	srv := cloudWatchServer("OK")
	defer srv.Close()

	c := Check{
		Fields: map[string]string{
			"alarm_name": "test-alarm",
			"endpoint":   srv.URL,
		},
	}
	drift, _, err := checkCloudWatchAlarm(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift with default expected_state=OK")
	}
}

func TestCheckCloudWatchAlarm_InvalidURL(t *testing.T) {
	c := Check{
		Fields: map[string]string{
			"alarm_name": "test-alarm",
			"endpoint":   "http://127.0.0.1:0",
		},
	}
	_, _, err := checkCloudWatchAlarm(c)
	if err == nil {
		t.Fatal("expected error for unreachable endpoint")
	}
}

func TestCheckCloudWatchAlarm_ViaChecker(t *testing.T) {
	srv := cloudWatchServer("INSUFFICIENT_DATA")
	defer srv.Close()

	chkr := New()
	c := Check{
		Name: "cw-test",
		Type: "cloudwatch_alarm",
		Fields: map[string]string{
			"alarm_name":     "test-alarm",
			"expected_state": "INSUFFICIENT_DATA",
			"endpoint":       srv.URL,
		},
	}
	drift, _, err := chkr.Run(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift")
	}
}
