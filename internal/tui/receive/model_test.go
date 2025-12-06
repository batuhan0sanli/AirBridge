package receive

import (
	"AirBridge/internal/crypto"
	"strings"
	"testing"
)

func TestInitialModel(t *testing.T) {
	// 1. Test with empty key (Default behavior)
	model1 := InitialModel(nil)
	if model1.step != StepGeneratingKey {
		t.Errorf("Expected step StepGeneratingKey for empty input, got %v", model1.step)
	}
	if model1.privateKey != nil {
		t.Error("Expected nil private key for empty input")
	}

	// 2. Test with Valid Private Key
	privKey, _, err := crypto.GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}
	privKeyPEM, err := crypto.ExportRSAPrivateKeyAsPEM(privKey)
	if err != nil {
		t.Fatalf("Failed to export private key: %v", err)
	}

	model2 := InitialModel(privKeyPEM)
	if model2.step != StepAwaitingPayload {
		t.Errorf("Expected step StepAwaitingPayload for valid key, got %v", model2.step)
	}
	if model2.privateKey == nil {
		t.Error("Expected private key to be set")
	}
	if model2.publicKey == nil {
		t.Error("Expected public key to be set")
	}
	if model2.encodedKey == "" {
		t.Error("Expected encoded public key to be set")
	}
	if model2.statusText != "" {
		t.Errorf("Expected empty status text, got %q", model2.statusText)
	}

	// 3. Test with Invalid Private Key
	invalidKeyPEM := []byte("-----BEGIN RSA PRIVATE KEY-----\nINVALID_DATA\n-----END RSA PRIVATE KEY-----")
	model3 := InitialModel(invalidKeyPEM)

	// Should fall back to generating key
	if model3.step != StepGeneratingKey {
		t.Errorf("Expected step StepGeneratingKey for invalid key, got %v", model3.step)
	}
	if model3.privateKey != nil {
		t.Error("Expected nil private key for invalid input")
	}
	// Check for warning message
	if !strings.Contains(model3.statusText, "Invalid private key provided") {
		t.Errorf("Expected warning in status text, got %q", model3.statusText)
	}
}
