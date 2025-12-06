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
	defer func() { _ = os.Remove(tmpfile.Name()) }() // clean up

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

	nonce, err := GenerateIV()
	if err != nil {
		t.Fatalf("GenerateIV failed: %v", err)
	}
	// GCM nonce is 12 bytes
	if len(nonce) != 12 {
		t.Errorf("Expected nonce length 12, got %d", len(nonce))
	}

	data := []byte("secret message")
	encrypted, err := EncryptDataAES(key, nonce, data)
	if err != nil {
		t.Fatalf("EncryptDataAES failed: %v", err)
	}

	// Decrypt to verify
	decrypted, err := DecryptDataAES(key, nonce, encrypted)
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

func TestGenerateAndEncodeRSAKeys(t *testing.T) {
	// Test GenerateRSAKeyPair
	privateKey, publicKey, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("GenerateRSAKeyPair failed: %v", err)
	}
	if privateKey == nil || publicKey == nil {
		t.Fatal("Generated keys are nil")
	}

	// Test EncodeRSAPublicKey
	encodedKey, err := EncodeRSAPublicKey(publicKey)
	if err != nil {
		t.Fatalf("EncodeRSAPublicKey failed: %v", err)
	}
	if encodedKey == "" {
		t.Error("Encoded key is empty")
	}

	// Verify we can decode it back
	decodedKey, err := DecodeRSAPublicKey(encodedKey)
	if err != nil {
		t.Fatalf("DecodeRSAPublicKey failed: %v", err)
	}

	if decodedKey.N.Cmp(publicKey.N) != 0 || decodedKey.E != publicKey.E {
		t.Error("Decoded key does not match original key")
	}
}

func TestDecryptAESKeyWithRSA(t *testing.T) {
	privateKey, publicKey, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	aesKey := []byte("12345678901234567890123456789012")
	encryptedKey, err := EncryptAESKeyWithRSA(publicKey, aesKey)
	if err != nil {
		t.Fatal(err)
	}

	decryptedKey, err := DecryptAESKeyWithRSA(privateKey, encryptedKey)
	if err != nil {
		t.Fatalf("DecryptAESKeyWithRSA failed: %v", err)
	}

	if string(decryptedKey) != string(aesKey) {
		t.Errorf("Expected key %s, got %s", aesKey, decryptedKey)
	}
}

func TestDecryptDataAES(t *testing.T) {
	key, _ := GenerateAESKey()
	nonce, _ := GenerateIV()
	data := []byte("test data for decryption wrapper")

	encrypted, err := EncryptDataAES(key, nonce, data)
	if err != nil {
		t.Fatal(err)
	}

	// Test the wrapper function
	decrypted, err := DecryptDataAES(key, nonce, encrypted)
	if err != nil {
		t.Fatalf("DecryptDataAES failed: %v", err)
	}

	if string(decrypted) != string(data) {
		t.Errorf("Expected %s, got %s", data, decrypted)
	}
}

func TestExportAndDecodePEMKeys(t *testing.T) {
	// 1. Generate RSA Key Pair
	privateKey, publicKey, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("GenerateRSAKeyPair failed: %v", err)
	}

	// 2. Test ExportRSAPublicKeyAsPEM
	pubPEM, err := ExportRSAPublicKeyAsPEM(publicKey)
	if err != nil {
		t.Fatalf("ExportRSAPublicKeyAsPEM failed: %v", err)
	}
	if len(pubPEM) == 0 {
		t.Error("Exported public key PEM is empty")
	}

	// 3. Test DecodeRSAPublicKey with raw PEM (string)
	// This verifies the fix where we handle non-base64 input
	pubPEMStr := string(pubPEM)
	decodedKey, err := DecodeRSAPublicKey(pubPEMStr)
	if err != nil {
		t.Fatalf("DecodeRSAPublicKey failed with raw PEM: %v", err)
	}

	if decodedKey.N.Cmp(publicKey.N) != 0 || decodedKey.E != publicKey.E {
		t.Error("Decoded key does not match original key")
	}

	// 4. Test ExportRSAPrivateKeyAsPEM
	privPEM, err := ExportRSAPrivateKeyAsPEM(privateKey)
	if err != nil {
		t.Fatalf("ExportRSAPrivateKeyAsPEM failed: %v", err)
	}
	if len(privPEM) == 0 {
		t.Error("Exported private key PEM is empty")
	}

	// Basic check for PEM header
	expectedHeader := "-----BEGIN RSA PRIVATE KEY-----"
	if string(privPEM[:len(expectedHeader)]) != expectedHeader {
		t.Errorf("Expected PEM header %q, got %q", expectedHeader, string(privPEM[:len(expectedHeader)]))
	}
}

func TestDecodeRSAPrivateKey(t *testing.T) {
	// 1. Generate keys
	privateKey, _, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("GenerateRSAKeyPair failed: %v", err)
	}

	// 2. Export Private Key to PEM
	privPEM, err := ExportRSAPrivateKeyAsPEM(privateKey)
	if err != nil {
		t.Fatalf("ExportRSAPrivateKeyAsPEM failed: %v", err)
	}

	// 3. Decode back
	decodedPrivKey, err := DecodeRSAPrivateKey(privPEM)
	if err != nil {
		t.Fatalf("DecodeRSAPrivateKey failed: %v", err)
	}

	// 4. Validate
	if decodedPrivKey.N.Cmp(privateKey.N) != 0 {
		t.Error("Decoded private key N does not match original key")
	}
	if decodedPrivKey.D.Cmp(privateKey.D) != 0 {
		t.Error("Decoded private key D does not match original key")
	}
}
