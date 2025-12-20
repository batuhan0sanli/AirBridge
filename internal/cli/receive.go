package cli

import (
	"AirBridge/internal/crypto"
	"AirBridge/pkg"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ProcessPayload decodes, decrypts and saves the file from the payload
func ProcessPayload(payloadStr string, privateKey *rsa.PrivateKey) (string, error) {
	// 1. Parse Base64 payload
	jsonPayloadBytes, err := base64.StdEncoding.DecodeString(payloadStr)
	if err != nil {
		return "", fmt.Errorf("invalid base64 payload: %v", err)
	}

	var payload pkg.SmallFilePayload
	if err := json.Unmarshal(jsonPayloadBytes, &payload); err != nil {
		return "", fmt.Errorf("invalid json payload: %v", err)
	}

	// 2. Decrypt AES Key
	encryptedAESKey, err := hex.DecodeString(payload.Key)
	if err != nil {
		return "", fmt.Errorf("invalid hex key: %v", err)
	}

	aesKey, err := crypto.DecryptAESKeyWithRSA(privateKey, encryptedAESKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt AES key: %v", err)
	}

	// 3. Decrypt Data
	nonce, err := hex.DecodeString(payload.Nonce)
	if err != nil {
		return "", fmt.Errorf("invalid hex nonce: %v", err)
	}

	encryptedData, err := hex.DecodeString(payload.Data)
	if err != nil {
		return "", fmt.Errorf("invalid hex data: %v", err)
	}

	decryptedData, err := crypto.DecryptDataAES(aesKey, nonce, encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt data: %v", err)
	}

	// 4. Save File
	safeFilename := filepath.Base(payload.Metadata.Name)
	err = os.WriteFile(safeFilename, decryptedData, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return safeFilename, nil
}

// RunReceive orchestrates the headless receive command
func RunReceive(payload string, privKeyPEM []byte, inputPayloadPath string, deletePayload bool) error {
	privKey, err := crypto.DecodeRSAPrivateKey(privKeyPEM)
	if err != nil {
		return fmt.Errorf("error decoding private key: %w", err)
	}

	filename, err := ProcessPayload(payload, privKey)
	if err != nil {
		return fmt.Errorf("error processing payload: %w", err)
	}

	fmt.Printf("File saved successfully: %s\n", filename)

	if deletePayload {
		err := os.Remove(inputPayloadPath)
		if err != nil {
			fmt.Printf("Warning: Failed to delete payload file: %v\n", err)
		} else {
			fmt.Println("Payload file deleted.")
		}
	}
	return nil
}
