package receive

import (
	"AirBridge/internal/crypto"
	"AirBridge/pkg"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestGenerateKeyCmd(t *testing.T) {
	cmd := generateKeyCmd()
	msg := cmd()

	keyMsg, ok := msg.(keyGeneratedMsg)
	if !ok {
		t.Fatalf("Expected keyGeneratedMsg, got %T", msg)
	}

	if keyMsg.privateKey == nil {
		t.Error("Private key is nil")
	}
	if keyMsg.publicKey == nil {
		t.Error("Public key is nil")
	}
	if keyMsg.encodedKey == "" {
		t.Error("Encoded key is empty")
	}
}

// Helper function to simulate encryption (copied from send/service.go logic for testing)
func encryptForTest(t *testing.T, data []byte, pubKey *rsa.PublicKey) string {
	aesKey, err := crypto.GenerateAESKey()
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	encryptedAESKey, err := crypto.EncryptAESKeyWithRSA(pubKey, aesKey)
	if err != nil {
		t.Fatalf("Failed to encrypt AES key: %v", err)
	}

	// Generate Nonce for AES-GCM
	nonce, err := crypto.GenerateIV()
	if err != nil {
		t.Fatalf("Failed to generate nonce: %v", err)
	}

	encryptedData, err := crypto.EncryptDataAES(aesKey, nonce, data)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	payload := pkg.SmallFilePayload{
		Key:      fmt.Sprintf("%x", encryptedAESKey),
		Data:     fmt.Sprintf("%x", encryptedData),
		Nonce:    fmt.Sprintf("%x", nonce),
		Metadata: pkg.FileMetadata{Name: "test_decrypted.txt", Size: int64(len(data)), Hash: "dummy"},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	return base64.StdEncoding.EncodeToString(jsonPayload)
}

func TestDecryptAndSaveCmd(t *testing.T) {
	// 1. Generate keys
	privKey, pubKey, err := crypto.GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	// 2. Prepare encrypted payload
	originalData := []byte("secret message for receive test")
	payloadStr := encryptForTest(t, originalData, pubKey)

	// 3. Run the command
	cmd := decryptAndSaveCmd(payloadStr, privKey)
	msg := cmd()

	// 4. Check result
	if errMsg, ok := msg.(errMsg); ok {
		t.Fatalf("Command returned error: %v", errMsg.error)
	}

	if _, ok := msg.(fileDecryptedMsg); !ok {
		t.Fatalf("Expected fileDecryptedMsg, got %T", msg)
	}

	// 5. Verify file content
	savedContent, err := os.ReadFile("test_decrypted.txt")
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if string(savedContent) != string(originalData) {
		t.Errorf("Expected content %q, got %q", string(originalData), string(savedContent))
	}

	// Cleanup
	if err := os.Remove("test_decrypted.txt"); err != nil {
		t.Logf("Failed to remove test file: %v", err)
	}
}

func TestDecryptAndSaveCmd_InvalidPayload(t *testing.T) {
	privKey, _, _ := crypto.GenerateRSAKeyPair()
	cmd := decryptAndSaveCmd("invalid_base64", privKey)
	msg := cmd()

	if _, ok := msg.(errMsg); !ok {
		t.Error("Expected errMsg for invalid payload")
	}
}

func TestDecryptAndSaveCmd_PathTraversal(t *testing.T) {
	// 1. Generate keys
	privKey, pubKey, err := crypto.GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	// 2. Prepare encrypted payload with malicious filename
	originalData := []byte("safe content")

	// Manually construct payload to inject malicious name
	aesKey, err := crypto.GenerateAESKey()
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	encryptedAESKey, err := crypto.EncryptAESKeyWithRSA(pubKey, aesKey)
	if err != nil {
		t.Fatalf("Failed to encrypt AES key: %v", err)
	}

	nonce, err := crypto.GenerateIV()
	if err != nil {
		t.Fatalf("Failed to generate nonce: %v", err)
	}

	encryptedData, err := crypto.EncryptDataAES(aesKey, nonce, originalData)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Malicious filename attempting traversal
	maliciousName := "../traversal_test.txt"

	payload := pkg.SmallFilePayload{
		Key:      fmt.Sprintf("%x", encryptedAESKey),
		Data:     fmt.Sprintf("%x", encryptedData),
		Nonce:    fmt.Sprintf("%x", nonce),
		Metadata: pkg.FileMetadata{Name: maliciousName, Size: int64(len(originalData)), Hash: "dummy"},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	payloadStr := base64.StdEncoding.EncodeToString(jsonPayload)

	// 3. Run the command
	cmd := decryptAndSaveCmd(payloadStr, privKey)
	msg := cmd()

	// 4. Check result
	if errMsg, ok := msg.(errMsg); ok {
		t.Fatalf("Command returned error: %v", errMsg.error)
	}

	// 5. Verify file location
	// It should NOT be in the parent directory
	if _, err := os.Stat("../traversal_test.txt"); err == nil {
		_ = os.Remove("../traversal_test.txt")
		t.Fatal("Security regression: File was written to ../traversal_test.txt")
	}

	// It SHOULD be in the current directory (sanitized)
	expectedFile := "traversal_test.txt"
	if _, err := os.Stat(expectedFile); err != nil {
		t.Fatalf("File was not written to expected sanitized path %s: %v", expectedFile, err)
	}

	// Cleanup
	_ = os.Remove(expectedFile)
}
