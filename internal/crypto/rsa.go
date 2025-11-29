package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// DecodeRSAPublicKey decodes a base64 encoded PEM string into an RSA public key.
func DecodeRSAPublicKey(pubKeyStr string) (*rsa.PublicKey, error) {
	// Base64 decode
	pemBytes, err := base64.StdEncoding.DecodeString(pubKeyStr)
	if err != nil {
		return nil, fmt.Errorf("could not decode base64 public key: %v", err)
	}

	// Decode PEM block
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, fmt.Errorf("could not decode PEM block")
	}

	// Parse PKIX public key
	genericPublicKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// Assert RSA public key
	rsaPublicKey, ok := genericPublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("the provided key is not an RSA public key")
	}

	return rsaPublicKey, nil
}

// EncryptAESKeyWithRSA encrypts an AES key using RSA-OAEP with SHA256.
func EncryptAESKeyWithRSA(publicKey *rsa.PublicKey, aesKey []byte) ([]byte, error) {
	encryptedAESKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, aesKey, nil)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt symmetric key with public key: %v", err)
	}
	return encryptedAESKey, nil
}
