package checker

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkSSLExpiry connects to a host and checks whether the TLS certificate
// expires within a threshold number of days. Fields used:
//   - "host": hostname:port (e.g. "example.com:443")
//   - "days": minimum days until expiry before drift is reported (default 14)
func checkSSLExpiry(check config.Check) (bool, string, error) {
	host, ok := check.Fields["host"]
	if !ok || host == "" {
		return false, "", fmt.Errorf("ssl_expiry check %q: missing required field 'host'", check.Name)
	}

	thresholdDays := 14
	if daysStr, ok := check.Fields["days"]; ok && daysStr != "" {
		_, err := fmt.Sscanf(daysStr, "%d", &thresholdDays)
		if err != nil {
			return false, "", fmt.Errorf("ssl_expiry check %q: invalid 'days' value: %s", check.Name, daysStr)
		}
		if thresholdDays < 0 {
			return false, "", fmt.Errorf("ssl_expiry check %q: 'days' must be non-negative, got %d", check.Name, thresholdDays)
		}
	}

	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: 5 * time.Second},
		"tcp",
		host,
		&tls.Config{InsecureSkipVerify: false},
	)
	if err != nil {
		return false, "", fmt.Errorf("ssl_expiry check %q: TLS dial failed: %w", check.Name, err)
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return false, "", fmt.Errorf("ssl_expiry check %q: no peer certificates found", check.Name)
	}

	expiry := certs[0].NotAfter
	daysRemaining := int(time.Until(expiry).Hours() / 24)
	detail := fmt.Sprintf("cert for %s expires in %d days (threshold: %d)", host, daysRemaining, thresholdDays)

	if daysRemaining < thresholdDays {
		return true, detail, nil
	}
	return false, detail, nil
}
