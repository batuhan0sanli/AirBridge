package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

// GenerateAESKey generates a random 32-byte AES key.
func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("could not generate symmetric key: %v", err)
	}
	return key, nil
}

// GenerateIV generates a random 12-byte Nonce for AES-GCM.
func GenerateIV() ([]byte, error) {
	// GCM standard nonce size is 12 bytes
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("could not generate nonce: %v", err)
	}
	return nonce, nil
}

// EncryptDataAES encrypts data using AES-GCM.
func EncryptDataAES(key, nonce, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create AES cipher: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("could not create GCM: %v", err)
	}

	// Seal encrypts and authenticates the data.
	// The nonce is passed as the first argument to Seal, but usually handling it separately is better for storage.
	// Here we just return the ciphertext (which includes the tag appended to it).
	ciphertext := aesGCM.Seal(nil, nonce, data, nil)

	return ciphertext, nil
}

// DecryptDataAES decrypts data using AES-GCM.
func DecryptDataAES(key, nonce, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create AES cipher: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("could not create GCM: %v", err)
	}

	plaintext, err := aesGCM.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt/authenticate data: %v", err)
	}

	return plaintext, nil
}
