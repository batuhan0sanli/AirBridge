package crypto

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// CalculateFileHash computes the SHA256 hash of a file.
// It resets the file pointer to the beginning before and after reading.
func CalculateFileHash(file *os.File) (string, error) {
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	fileHash := hasher.Sum(nil)

	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", fileHash), nil
}
