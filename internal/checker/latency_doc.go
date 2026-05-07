// Package checker provides the http_latency check type.
//
// # HTTP Latency Check
//
// The http_latency check probes an HTTP endpoint and reports drift when the
// round-trip response time exceeds a configured threshold.
//
// ## Configuration fields
//
//	url      (required) – Full URL to probe, e.g. https://api.example.com/health
//	max_ms   (required) – Maximum acceptable latency in milliseconds (positive integer)
//	method   (optional) – HTTP method; defaults to GET
//
// ## Example
//
//	checks:
//	  - name: api-latency
//	    type: http_latency
//	    fields:
//	      url: https://api.example.com/health
//	      max_ms: "300"
//	      method: GET
package checker
