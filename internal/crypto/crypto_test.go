package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"testing"
)

func TestCalculateFileHash(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	content := []byte("hello world")
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}

	// Calculate hash
	hash, err := CalculateFileHash(tmpfile)
	if err != nil {
		t.Fatalf("CalculateFileHash failed: %v", err)
	}

	expectedHash := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	if hash != expectedHash {
		t.Errorf("Expected hash %s, got %s", expectedHash, hash)
	}
}

func TestAESEncryption(t *testing.T) {
	key, err := GenerateAESKey()
	if err != nil {
		t.Fatalf("GenerateAESKey failed: %v", err)
	}
	if len(key) != 32 {
		t.Errorf("Expected key length 32, got %d", len(key))
	}

	iv, err := GenerateIV()
	if err != nil {
		t.Fatalf("GenerateIV failed: %v", err)
	}
	if len(iv) != 16 {
		t.Errorf("Expected IV length 16, got %d", len(iv))
	}

	data := []byte("secret message")
	encrypted, err := EncryptDataAES(key, iv, data)
	if err != nil {
		t.Fatalf("EncryptDataAES failed: %v", err)
	}

	// Decrypt to verify (using the same function since CTR is symmetric)
	decrypted, err := EncryptDataAES(key, iv, encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != string(data) {
		t.Errorf("Expected %s, got %s", data, decrypted)
	}
}

func TestRSAEncryption(t *testing.T) {
	// Generate RSA key pair for testing
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	publicKey := &privateKey.PublicKey

	// Encode public key to PEM format to test DecodeRSAPublicKey
	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatal(err)
	}
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	pubStr := base64.StdEncoding.EncodeToString(pubBytes)

	// Test DecodeRSAPublicKey
	decodedKey, err := DecodeRSAPublicKey(pubStr)
	if err != nil {
		t.Fatalf("DecodeRSAPublicKey failed: %v", err)
	}

	if decodedKey.N.Cmp(publicKey.N) != 0 || decodedKey.E != publicKey.E {
		t.Error("Decoded key does not match original key")
	}

	// Test EncryptAESKeyWithRSA
	aesKey := []byte("12345678901234567890123456789012")
	encryptedKey, err := EncryptAESKeyWithRSA(decodedKey, aesKey)
	if err != nil {
		t.Fatalf("EncryptAESKeyWithRSA failed: %v", err)
	}

	if len(encryptedKey) == 0 {
		t.Error("Encrypted key is empty")
	}
}
