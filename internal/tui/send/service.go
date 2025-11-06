package send

import (
	"AirBridge/pkg"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func getMetadata(file *os.File) (pkg.FileMetadata, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return pkg.FileMetadata{}, err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return pkg.FileMetadata{}, err
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return pkg.FileMetadata{}, err
	}
	fileHash := hasher.Sum(nil)

	return pkg.FileMetadata{
		Name: fileInfo.Name(),
		Size: fileInfo.Size(),
		Hash: fmt.Sprintf("%x", fileHash),
	}, nil
}

func decodePublicKey(pubKeyStr string) (*rsa.PublicKey, error) {
	// Base64 decode et
	pemBytes, err := base64.StdEncoding.DecodeString(pubKeyStr)
	if err != nil {
		return nil, fmt.Errorf("could not decode base64 public key: %v", err)
	}

	// 1. PEM bloğunu metinden çöz
	pemBlock, _ := pem.Decode([]byte(pemBytes))
	if pemBlock == nil {
		return nil, fmt.Errorf("could not decode PEM block")
	}

	// 2. PEM bloğunu x509 (PKIX) formatından Go'nun anlayacağı public key formatına çevir
	genericPublicKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// 3. Anahtarın bir RSA public key olduğunu doğrula
	rsaPublicKey, ok := genericPublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("the provided key is not an RSA public key")
	}

	fmt.Println("Public key successfully parsed.")
	return rsaPublicKey, nil
}

// message types for async workflow (split steps)
type fileOpenedMsg struct{ file *os.File }
type metadataExtractedMsg struct{ metadata pkg.FileMetadata }
type publicKeyProcessedMsg struct{ publicKey *rsa.PublicKey }
type errMsg struct{ error }

// openFileCmd opens the file asynchronously
func openFileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		// debugging amaçlı bekleme komut fonksiyonu içinde olmalı ki UI bloklanmasın
		time.Sleep(1 * time.Second)
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
		// debugging amaçlı bekleme komut fonksiyonu içinde olmalı ki UI bloklanmasın
		time.Sleep(1 * time.Second)
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

func processPublicKeyCmd(rawPublicKey string) tea.Cmd {
	return func() tea.Msg {
		// debugging amaçlı bekleme komut fonksiyonu içinde olmalı ki UI bloklanmasın
		time.Sleep(1 * time.Second)

		pubKey, err := decodePublicKey(rawPublicKey)
		if err != nil {
			return errMsg{err}
		}
		return publicKeyProcessedMsg{publicKey: pubKey}
	}
}
