package encryptor

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

const divider = "_"
const chunksKey = "."

type Encryptor struct {
	privateKey       *rsa.PrivateKey
	publicKey        *rsa.PublicKey
	channelPublicKey *rsa.PublicKey
}

func GenerateKeyPair() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		log.Fatalf("failed to generate private key: %v", err)
	}

	return privateKey
}

func NewEncryptor(privateKey *rsa.PrivateKey) Encryptor {
	return Encryptor{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}
}

func (b *Encryptor) GetPublicKey() *rsa.PublicKey {
	return b.publicKey
}

func (b *Encryptor) SetChannelPublicKey(channelPublicKey *rsa.PublicKey) {
	b.channelPublicKey = channelPublicKey
}
func (b *Encryptor) Encrypt(value string) (string, error) {
	hash := sha256.Sum256([]byte(value))
	signature, err := rsa.SignPKCS1v15(rand.Reader, b.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	maxChunkSize := b.channelPublicKey.Size() - 11
	chunks := splitBytes([]byte(value), maxChunkSize)

	var encChunks []string
	for _, chunk := range chunks {
		enc, err := rsa.EncryptPKCS1v15(rand.Reader, b.channelPublicKey, chunk)
		if err != nil {
			return "", err
		}
		encChunks = append(encChunks, base64.StdEncoding.EncodeToString(enc))
	}

	result := fmt.Sprintf("%s%s%s",
		strings.Join(encChunks, chunksKey), // chunk separator
		divider,
		base64.StdEncoding.EncodeToString(signature),
	)

	return base64.StdEncoding.EncodeToString([]byte(result)), nil
}

func (b *Encryptor) Decrypt(value string) (string, error) {
	decodedValue, err := base64.StdEncoding.DecodeString(value)
	parts := strings.Split(string(decodedValue), divider)

	if len(parts) != 2 {
		return "", errors.New("invalid encrypted value")
	}
	chunkEncoded := strings.Split(parts[0], chunksKey)
	signature, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	var plainBuf bytes.Buffer
	for _, encChunk := range chunkEncoded {
		if len(encChunk) == 0 {
			continue
		}
		cipherChunk, err := base64.StdEncoding.DecodeString(encChunk)
		if err != nil {
			return "", err
		}
		plainChunk, err := rsa.DecryptPKCS1v15(rand.Reader, b.privateKey, cipherChunk)
		if err != nil {
			return "", err
		}
		plainBuf.Write(plainChunk)
	}
	message := plainBuf.String()

	// Verify the signature on the full message
	hash2 := sha256.Sum256([]byte(message))
	err = rsa.VerifyPKCS1v15(b.channelPublicKey, crypto.SHA256, hash2[:], signature)
	if err != nil {
		return "", errors.New("signature verification failed")
	}

	return message, nil
}

func splitBytes(data []byte, chunkSize int) [][]byte {
	var chunks [][]byte
	for len(data) > 0 {
		n := chunkSize
		if len(data) < chunkSize {
			n = len(data)
		}
		chunks = append(chunks, data[:n])
		data = data[n:]
	}
	return chunks
}
