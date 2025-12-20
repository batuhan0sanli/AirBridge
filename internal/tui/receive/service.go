package receive

import (
	"AirBridge/internal/cli"
	"AirBridge/internal/crypto"
	"crypto/rsa"

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
		_, err := cli.ProcessPayload(payloadStr, privateKey)
		if err != nil {
			return errMsg{err}
		}
		return fileDecryptedMsg{}
	}
}
