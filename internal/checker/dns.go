package checker

import (
	"fmt"
	"net"
	"strings"

	"github.com/user/driftwatch/internal/config"
)

// checkDNSResolve checks that a hostname resolves to an expected IP address.
// Config fields:
//   - host:     the hostname to resolve (required)
//   - expected: comma-separated list of expected IPs (required)
func checkDNSResolve(check config.Check) (bool, string, error) {
	host, ok := check.Config["host"]
	if !ok || strings.TrimSpace(host) == "" {
		return false, "", fmt.Errorf("dns check %q: missing or empty 'host' field", check.Name)
	}

	expected, ok := check.Config["expected"]
	if !ok || strings.TrimSpace(expected) == "" {
		return false, "", fmt.Errorf("dns check %q: missing or empty 'expected' field", check.Name)
	}

	addrs, err := net.LookupHost(host)
	if err != nil {
		return false, "", fmt.Errorf("dns check %q: lookup failed: %w", check.Name, err)
	}

	expectedSet := make(map[string]struct{})
	for _, e := range strings.Split(expected, ",") {
		expectedSet[strings.TrimSpace(e)] = struct{}{}
	}

	for _, addr := range addrs {
		if _, found := expectedSet[addr]; found {
			return false, "", nil // no drift
		}
	}

	return true,
		fmt.Sprintf("DNS drift: host %q resolved to %v, expected one of [%s]", host, addrs, expected),
		nil
}
