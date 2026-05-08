// Package checker provides infrastructure drift detection checks.
//
// # certpin check
//
// The certpin check connects to a TLS endpoint and verifies that the leaf
// certificate's SHA-256 fingerprint matches a known-good pin. This guards
// against unexpected certificate rotations or supply-chain substitution.
//
// Required params:
//   - host   (string) – hostname or IP to connect to
//   - pin    (string) – expected SHA-256 fingerprint in lowercase hex
//
// Optional params:
//   - port   (string) – TCP port to connect on (default: "443")
//
// Example YAML:
//
//	- name: prod-api-cert-pin
//	  type: certpin
//	  params:
//	    host: api.example.com
//	    pin: "a3b2c1..."   # sha256 hex fingerprint of leaf cert
//	    port: "8443"
package checker
