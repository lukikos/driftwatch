// Package checker provides the jsonfield check type.
//
// # jsonfield
//
// Fetches a JSON HTTP endpoint and asserts that a specific field
// equals an expected value. Useful for monitoring API health responses,
// feature flags, or any JSON-emitting service endpoint.
//
// ## Configuration fields
//
//	- url      (required): The HTTP or HTTPS URL to fetch.
//	- field    (required): Dot-separated path to the target field (e.g. "status" or "db.connected").
//	- expected (required): The expected string representation of the field value.
//
// ## Example
//
//	checks:
//	  - name: api-health-status
//	    type: jsonfield
//	    fields:
//	      url: https://api.example.com/health
//	      field: status
//	      expected: ok
//
//	  - name: db-connection-check
//	    type: jsonfield
//	    fields:
//	      url: https://api.example.com/health
//	      field: db.connected
//	      expected: "true"
package checker
