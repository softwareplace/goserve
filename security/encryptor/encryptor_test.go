package encryptor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Client struct {
	Encryptor
}

func NewClient() *Client {
	return &Client{
		Encryptor: NewEncryptor(GenerateKeyPair()),
	}
}

func TestEncryptorValidator(t *testing.T) {
	t.Run("should encrypt and decrypt message between 2 clients", func(t *testing.T) {
		client1 := NewClient()
		client2 := NewClient()
		client1.SetChannelPublicKey(client2.GetPublicKey())
		client2.SetChannelPublicKey(client1.GetPublicKey())

		message := "Hello encryptor"

		encryptedMessage, err := client1.Encrypt(message)
		require.NoError(t, err)
		require.NotEqual(t, message, encryptedMessage)

		decryptedMessage, err := client2.Decrypt(encryptedMessage)

		require.NoError(t, err)
		require.Equal(t, message, decryptedMessage)
	})

	t.Run("should encrypt and decrypt message between from the same client", func(t *testing.T) {
		client1 := NewClient()
		client1.SetChannelPublicKey(client1.GetPublicKey())

		message := "Hello encryptor"

		encryptedMessage, err := client1.Encrypt(message)
		require.NoError(t, err)
		require.NotEqual(t, message, encryptedMessage)

		decryptedMessage, err := client1.Decrypt(encryptedMessage)

		require.NoError(t, err)
		require.Equal(t, message, decryptedMessage)
	})

}
