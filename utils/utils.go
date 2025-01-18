package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

func Encrypt(value string, secret []byte) (string, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	// Generate a random IV
	iv := make([]byte, block.BlockSize())
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", fmt.Errorf("failed to generate IV: %w", err)
	}

	// Encrypt the data
	keyBytes := []byte(value)
	encrypted := make([]byte, len(keyBytes))
	mode := cipher.NewCFBEncrypter(block, iv)
	mode.XORKeyStream(encrypted, keyBytes)

	// Prepend the IV to the encrypted data and encode as Base64
	result := append(iv, encrypted...)
	return base64.StdEncoding.EncodeToString(result), nil
}

func Decrypt(encrypted string, secret []byte) (string, error) {
	if len(strings.TrimSpace(encrypted)) == 0 {
		return "", fmt.Errorf("encrypted value must not be empty")
	}
	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	// Decode the Base64-encoded encrypted data
	encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode Base64 data: %w", err)
	}

	// Extract the IV and encrypted data
	iv := encryptedBytes[:block.BlockSize()]
	encryptedData := encryptedBytes[block.BlockSize():]

	// Decrypt the data
	decrypted := make([]byte, len(encryptedData))
	mode := cipher.NewCFBDecrypter(block, iv)
	mode.XORKeyStream(decrypted, encryptedData)
	return string(decrypted), nil
}
