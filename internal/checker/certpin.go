package checker

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"crypto/sha256"
	"fmt"
	"net"
	"time"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkCertPin connects to a TLS host and verifies that the leaf certificate's
// SHA-256 fingerprint matches the expected pin. Drift is reported when the
// fingerprint does not match or the connection cannot be established.
func checkCertPin(check config.Check) (bool, string, error) {
	host, ok := check.Params["host"].(string)
	if !ok || host == "" {
		return false, "", fmt.Errorf("certpin: missing or empty 'host' param")
	}

	expectedPin, ok := check.Params["pin"].(string)
	if !ok || expectedPin == "" {
		return false, "", fmt.Errorf("certpin: missing or empty 'pin' param (SHA-256 hex fingerprint)")
	}

	port := "443"
	if p, ok := check.Params["port"].(string); ok && p != "" {
		port = p
	}

	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: 10 * time.Second},
		"tcp",
		net.JoinHostPort(host, port),
		&tls.Config{InsecureSkipVerify: true}, //nolint:gosec // intentional: we verify via pin
	)
	if err != nil {
		return false, "", fmt.Errorf("certpin: failed to connect to %s:%s: %w", host, port, err)
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return false, "", fmt.Errorf("certpin: no certificates returned by %s:%s", host, port)
	}

	actualPin := certFingerprint(certs[0])
	if actualPin != expectedPin {
		return true, fmt.Sprintf("certificate pin mismatch: got %s, expected %s", actualPin, expectedPin), nil
	}

	return false, "certificate pin matches expected value", nil
}

func certFingerprint(cert *x509.Certificate) string {
	sum := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(sum[:])
}
