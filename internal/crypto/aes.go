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

// GenerateIV generates a random 16-byte Initialization Vector.
func GenerateIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("could not generate IV: %v", err)
	}
	return iv, nil
}

// EncryptDataAES encrypts data using AES-CTR.
func EncryptDataAES(key, iv, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create AES cipher: %v", err)
	}

	stream := cipher.NewCTR(block, iv)
	encryptedData := make([]byte, len(data))
	stream.XORKeyStream(encryptedData, data)

	return encryptedData, nil
}

// DecryptDataAES decrypts data using AES-CTR.
// Since CTR mode is symmetric, encryption and decryption are the same operation.
func DecryptDataAES(key, iv, data []byte) ([]byte, error) {
	return EncryptDataAES(key, iv, data)
}
