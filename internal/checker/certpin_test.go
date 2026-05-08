package checker

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/config"
)

func makeSelfSignedCert(t *testing.T) (tls.Certificate, string) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create cert: %v", err)
	}

	sum := sha256.Sum256(certDER)
	fingerprint := hex.EncodeToString(sum[:])

	cert, err := tls.X509KeyPair(
		pemEncode("CERTIFICATE", certDER),
		pemEncodeKey(key),
	)
	if err != nil {
		t.Fatalf("x509 key pair: %v", err)
	}

	return cert, fingerprint
}

func startTLSPinServer(t *testing.T, cert tls.Certificate) string {
	t.Helper()

	ln, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{Certificates: []tls.Certificate{cert}})
	if err != nil {
		t.Fatalf("tls listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return ln.Addr().String()
}

func TestCheckCertPin_NoDrift(t *testing.T) {
	cert, pin := makeSelfSignedCert(t)
	addr := startTLSPinServer(t, cert)
	host, port, _ := net.SplitHostPort(addr)

	check := config.Check{
		Name: "pin-ok",
		Type: "certpin",
		Params: map[string]interface{}{"host": host, "port": port, "pin": pin},
	}

	drift, msg, err := checkCertPin(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Errorf("expected no drift, got: %s", msg)
	}
}

func TestCheckCertPin_Drift(t *testing.T) {
	cert, _ := makeSelfSignedCert(t)
	addr := startTLSPinServer(t, cert)
	host, port, _ := net.SplitHostPort(addr)

	check := config.Check{
		Name: "pin-mismatch",
		Type: "certpin",
		Params: map[string]interface{}{"host": host, "port": port, "pin": "deadbeef"},
	}

	drift, msg, err := checkCertPin(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drift {
		t.Errorf("expected drift but got none; msg: %s", msg)
	}
}

func TestCheckCertPin_MissingHost(t *testing.T) {
	check := config.Check{
		Name:   "no-host",
		Type:   "certpin",
		Params: map[string]interface{}{"pin": "abc123"},
	}
	_, _, err := checkCertPin(check)
	if err == nil {
		t.Error("expected error for missing host")
	}
}

func TestCheckCertPin_MissingPin(t *testing.T) {
	check := config.Check{
		Name:   "no-pin",
		Type:   "certpin",
		Params: map[string]interface{}{"host": "example.com"},
	}
	_, _, err := checkCertPin(check)
	if err == nil {
		t.Error("expected error for missing pin")
	}
}

func TestCheckCertPin_ViaChecker(t *testing.T) {
	cert, pin := makeSelfSignedCert(t)
	addr := startTLSPinServer(t, cert)
	host, port, _ := net.SplitHostPort(addr)

	c := New()
	check := config.Check{
		Name:   "via-checker",
		Type:   "certpin",
		Params: map[string]interface{}{"host": host, "port": port, "pin": pin},
	}

	drift, _, err := c.Run(check)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drift {
		t.Error("expected no drift via checker dispatch")
	}
}
