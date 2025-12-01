package send

import (
	"AirBridge/internal/crypto"
	"AirBridge/pkg"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func getMetadata(file *os.File) (pkg.FileMetadata, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return pkg.FileMetadata{}, err
	}

	fileHash, err := crypto.CalculateFileHash(file)
	if err != nil {
		return pkg.FileMetadata{}, err
	}

	return pkg.FileMetadata{
		Name: fileInfo.Name(),
		Size: fileInfo.Size(),
		Hash: fileHash,
	}, nil
}

func encryptFile(file *os.File, metadata pkg.FileMetadata, publicKey *rsa.PublicKey) (string, error) {
	// 5. Encryption process (Generate random key for AES-256)
	aesKey, err := crypto.GenerateAESKey()
	if err != nil {
		return "", fmt.Errorf("Error: Could not generate symmetric key: %v\n", err)
	}

	// 6. Encrypt AES key with recipient's public key
	// Encrypt the AES key with the parsed rsaPublicKey.
	// OAEP is a modern and secure padding standard.
	encryptedAESKey, err := crypto.EncryptAESKeyWithRSA(publicKey, aesKey)
	if err != nil {
		return "", fmt.Errorf("Error: Could not encrypt symmetric key with public key: %v\n", err)
	}

	// Generate IV (Initialization Vector)
	iv, err := crypto.GenerateIV()
	if err != nil {
		return "", fmt.Errorf("Error: Could not generate IV: %v\n", err)
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Error reading small file into memory: %v\n", err)
	}

	// Encrypt with CTR stream
	encryptedData, err := crypto.EncryptDataAES(aesKey, iv, fileBytes)
	if err != nil {
		return "", fmt.Errorf("Error: Could not encrypt data: %v\n", err)
	}

	// Make Payload
	payload := pkg.SmallFilePayload{
		Key:      fmt.Sprintf("%x", encryptedAESKey),
		Data:     fmt.Sprintf("%x", encryptedData),
		IV:       fmt.Sprintf("%x", iv),
		Metadata: metadata,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("Error: Could not marshal JSON payload: %v\n", err)
	}

	encodedPayload := base64.StdEncoding.EncodeToString(jsonPayload)
	return encodedPayload, nil
}

// message types for async workflow (split steps)
type fileOpenedMsg struct{ file *os.File }
type metadataExtractedMsg struct{ metadata pkg.FileMetadata }

type smallFilePayloadMsg struct{ payload string }
type errMsg struct{ error }

// openFileCmd opens the file asynchronously
func openFileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		if path == "" {
			return errMsg{fmt.Errorf("empty file path")}
		}
		file, err := os.Open(path)
		if err != nil {
			return errMsg{err}
		}
		return fileOpenedMsg{file: file}
	}
}

// extractMetadataCmd extracts metadata asynchronously using an already opened file
func extractMetadataCmd(file *os.File) tea.Cmd {
	return func() tea.Msg {
		if file == nil {
			return errMsg{fmt.Errorf("nil file")}
		}
		metadata, err := getMetadata(file)
		if err != nil {
			return errMsg{err}
		}
		return metadataExtractedMsg{metadata: metadata}
	}
}

func processPublicKeyCmd(rawPublicKey string, file *os.File, metadata pkg.FileMetadata) tea.Cmd {
	return func() tea.Msg {
		pubKey, err := crypto.DecodeRSAPublicKey(rawPublicKey)
		if err != nil {
			return errMsg{err}
		}

		encryptedPayload, err := encryptFile(file, metadata, pubKey)
		if err != nil {
			return errMsg{err}
		}
		return smallFilePayloadMsg{payload: encryptedPayload}
	}
}
