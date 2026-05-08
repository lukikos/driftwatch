package checker

import (
	"crypto/ecdsa"
	"encoding/pem"
)

// pemEncode wraps DER bytes in a PEM block with the given type.
func pemEncode(blockType string, der []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: blockType, Bytes: der})
}

// pemEncodeKey marshals an ECDSA private key to PKCS#8 PEM.
func pemEncodeKey(key *ecdsa.PrivateKey) []byte {
	import_x509 := func() []byte {
		// Use MarshalECPrivateKey for simplicity in tests.
		import_crypto_x509, _ := func() ([]byte, error) {
			// inline to avoid circular import issues in test helpers
			return marshalECKey(key)
		}()
		return import_crypto_x509
	}
	return import_x509()
}

func marshalECKey(key *ecdsa.PrivateKey) ([]byte, error) {
	import (
		"crypto/x509"
	)
	der, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: der}), nil
}
