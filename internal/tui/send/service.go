package send

import (
	"AirBridge/internal/cli"
	"AirBridge/internal/crypto"
	"AirBridge/pkg"
	"crypto/rsa"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func getMetadata(file *os.File) (pkg.FileMetadata, error) {
	return cli.GetFileMetadata(file)
}

func encryptFile(file *os.File, metadata pkg.FileMetadata, publicKey *rsa.PublicKey) (string, error) {
	return cli.EncryptFile(file, metadata, publicKey)
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
