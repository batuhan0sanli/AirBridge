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

// GenerateRSAKeyPair generates a new RSA key pair of 2048 bits.
func GenerateRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("could not generate RSA key pair: %v", err)
	}
	return privateKey, &privateKey.PublicKey, nil
}

// EncodeRSAPublicKey encodes an RSA public key to a base64 encoded PEM string.
func EncodeRSAPublicKey(publicKey *rsa.PublicKey) (string, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("could not marshal public key: %v", err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})

	return base64.StdEncoding.EncodeToString(pubBytes), nil
}

// DecryptAESKeyWithRSA decrypts an AES key using RSA-OAEP with SHA256.
func DecryptAESKeyWithRSA(privateKey *rsa.PrivateKey, encryptedAESKey []byte) ([]byte, error) {
	decryptedAESKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedAESKey, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt symmetric key with private key: %v", err)
	}
	return decryptedAESKey, nil
}
