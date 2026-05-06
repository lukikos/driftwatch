package checker

import (
	"fmt"
	"net"
	"time"
)

// checkPortOpen checks whether a TCP port is open on a given host.
// Config fields:
//   - "host": hostname or IP (default: "localhost")
//   - "port": port number (required)
func checkPortOpen(fields map[string]string) (bool, string, error) {
	port, ok := fields["port"]
	if !ok || port == "" {
		return false, "", fmt.Errorf("check_port: missing required field 'port'")
	}

	host := fields["host"]
	if host == "" {
		host = "localhost"
	}

	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		// Port is not open — this is drift if we expect it to be open.
		return true, fmt.Sprintf("port %s on %s is not open: %v", port, host, err), nil
	}
	conn.Close()
	return false, "", nil
}
