package receive

import (
	"AirBridge/internal/crypto"
	"AirBridge/pkg"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type keyGeneratedMsg struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	encodedKey string
}

type fileDecryptedMsg struct{}

type errMsg struct{ error }

func generateKeyCmd() tea.Cmd {
	return func() tea.Msg {
		// Simulate delay for better UX
		time.Sleep(500 * time.Millisecond)

		privateKey, publicKey, err := crypto.GenerateRSAKeyPair()
		if err != nil {
			return errMsg{err}
		}

		encodedKey, err := crypto.EncodeRSAPublicKey(publicKey)
		if err != nil {
			return errMsg{err}
		}

		return keyGeneratedMsg{
			privateKey: privateKey,
			publicKey:  publicKey,
			encodedKey: encodedKey,
		}
	}
}

func decryptAndSaveCmd(payloadStr string, privateKey *rsa.PrivateKey) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)

		// 1. Parse Base64 payload
		jsonPayloadBytes, err := base64.StdEncoding.DecodeString(payloadStr)
		if err != nil {
			return errMsg{fmt.Errorf("invalid base64 payload: %v", err)}
		}

		var payload pkg.SmallFilePayload
		if err := json.Unmarshal(jsonPayloadBytes, &payload); err != nil {
			return errMsg{fmt.Errorf("invalid json payload: %v", err)}
		}

		// 2. Decrypt AES Key
		encryptedAESKey, err := hex.DecodeString(payload.Key)
		if err != nil {
			return errMsg{fmt.Errorf("invalid hex key: %v", err)}
		}

		aesKey, err := crypto.DecryptAESKeyWithRSA(privateKey, encryptedAESKey)
		if err != nil {
			return errMsg{fmt.Errorf("failed to decrypt AES key: %v", err)}
		}

		// 3. Decrypt Data
		iv, err := hex.DecodeString(payload.IV)
		if err != nil {
			return errMsg{fmt.Errorf("invalid hex iv: %v", err)}
		}

		encryptedData, err := hex.DecodeString(payload.Data)
		if err != nil {
			return errMsg{fmt.Errorf("invalid hex data: %v", err)}
		}

		decryptedData, err := crypto.DecryptDataAES(aesKey, iv, encryptedData)
		if err != nil {
			return errMsg{fmt.Errorf("failed to decrypt data: %v", err)}
		}

		// 4. Save File
		err = os.WriteFile(payload.Metadata.Name, decryptedData, 0644)
		if err != nil {
			return errMsg{fmt.Errorf("failed to save file: %v", err)}
		}

		return fileDecryptedMsg{}
	}
}
