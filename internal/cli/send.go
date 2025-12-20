package cli

import (
	"AirBridge/internal/crypto"
	"AirBridge/pkg"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// GetFileMetadata extracts metadata from the file
func GetFileMetadata(file *os.File) (pkg.FileMetadata, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return pkg.FileMetadata{}, err
	}

	fileHash, err := crypto.CalculateFileHash(file)
	if err != nil {
		return pkg.FileMetadata{}, err
	}

	// Reset file pointer after hash calculation
	_, err = file.Seek(0, 0)
	if err != nil {
		return pkg.FileMetadata{}, fmt.Errorf("failed to reset file pointer: %v", err)
	}

	return pkg.FileMetadata{
		Name: fileInfo.Name(),
		Size: fileInfo.Size(),
		Hash: fileHash,
	}, nil
}

// EncryptFile encrypts the file and returns the base64 encoded payload
func EncryptFile(file *os.File, metadata pkg.FileMetadata, publicKey *rsa.PublicKey) (string, error) {
	// 5. Encryption process (Generate random key for AES-256)
	aesKey, err := crypto.GenerateAESKey()
	if err != nil {
		return "", fmt.Errorf("could not generate symmetric key: %v", err)
	}

	// 6. Encrypt AES key with recipient's public key
	encryptedAESKey, err := crypto.EncryptAESKeyWithRSA(publicKey, aesKey)
	if err != nil {
		return "", fmt.Errorf("could not encrypt symmetric key with public key: %v", err)
	}

	// Generate Nonce (Number used once) for AES-GCM
	nonce, err := crypto.GenerateIV()
	if err != nil {
		return "", fmt.Errorf("could not generate nonce: %v", err)
	}

	// Ensure we read from start
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("failed to reset file pointer: %v", err)
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading file into memory: %v", err)
	}

	// Encrypt with GCM
	encryptedData, err := crypto.EncryptDataAES(aesKey, nonce, fileBytes)
	if err != nil {
		return "", fmt.Errorf("could not encrypt data: %v", err)
	}

	// Make Payload
	payload := pkg.SmallFilePayload{
		Key:      fmt.Sprintf("%x", encryptedAESKey),
		Data:     fmt.Sprintf("%x", encryptedData),
		Nonce:    fmt.Sprintf("%x", nonce),
		Metadata: metadata,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("could not marshal JSON payload: %v", err)
	}

	encodedPayload := base64.StdEncoding.EncodeToString(jsonPayload)
	return encodedPayload, nil
}

// RunSend orchestrates the headless send command
func RunSend(filePath string, pubKeyPEM string, outputFilePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	metadata, err := GetFileMetadata(file)
	if err != nil {
		return fmt.Errorf("error extracting metadata: %w", err)
	}

	pubKey, err := crypto.DecodeRSAPublicKey(pubKeyPEM)
	if err != nil {
		return fmt.Errorf("error decoding public key: %w", err)
	}

	payload, err := EncryptFile(file, metadata, pubKey)
	if err != nil {
		return fmt.Errorf("error encrypting file: %w", err)
	}

	outPath := outputFilePath
	if outPath == "" {
		outPath = "payload.abp"
	}

	err = os.WriteFile(outPath, []byte(payload), 0644)
	if err != nil {
		return fmt.Errorf("error saving payload: %w", err)
	}

	fmt.Printf("Payload saved to %s\n", outPath)
	return nil
}
