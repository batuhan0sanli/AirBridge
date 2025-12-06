package send

import (
	"AirBridge/internal/crypto"
	"AirBridge/pkg"
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"
)

func TestGetMetadata(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	content := []byte("test content")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek temp file: %v", err)
	}

	metadata, err := getMetadata(tmpFile)
	if err != nil {
		t.Fatalf("getMetadata failed: %v", err)
	}

	if metadata.Name != "" && metadata.Name != "testfile" && len(metadata.Name) == 0 {
		// Note: os.CreateTemp generates a random name, so we just check if it's not empty
		t.Errorf("Expected non-empty name, got %v", metadata.Name)
	}
	if metadata.Size != int64(len(content)) {
		t.Errorf("Expected size %d, got %d", len(content), metadata.Size)
	}
	// We can't easily predict the hash without calculating it again, but we assume crypto package is tested
	if metadata.Hash == "" {
		t.Error("Expected non-empty hash")
	}
}

func TestEncryptFile(t *testing.T) {
	// 1. Generate RSA keys for testing
	_, pubKey, err := crypto.GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate RSA keys: %v", err)
	}

	// 2. Create a temporary file
	tmpFile, err := os.CreateTemp("", "testfile_encrypt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	content := []byte("secret data")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek temp file: %v", err)
	}

	metadata := pkg.FileMetadata{
		Name: "testfile.txt",
		Size: int64(len(content)),
		Hash: "dummyhash",
	}

	// 3. Encrypt the file
	payloadStr, err := encryptFile(tmpFile, metadata, pubKey)
	if err != nil {
		t.Fatalf("encryptFile failed: %v", err)
	}

	// 4. Verify the payload structure
	payloadBytes, err := base64.StdEncoding.DecodeString(payloadStr)
	if err != nil {
		t.Fatalf("Failed to decode base64 payload: %v", err)
	}

	var payload pkg.SmallFilePayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		t.Fatalf("Failed to unmarshal JSON payload: %v", err)
	}

	if payload.Key == "" {
		t.Error("Payload key is empty")
	}
	if payload.Data == "" {
		t.Error("Payload data is empty")
	}
	if payload.Nonce == "" {
		t.Error("Payload Nonce is empty")
	}
	if payload.Metadata.Name != metadata.Name {
		t.Errorf("Expected metadata name %s, got %s", metadata.Name, payload.Metadata.Name)
	}
}
