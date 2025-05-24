package encryptor

import "fmt"

type serviceImpl struct {
	ApiSecretAuthorization []byte
}

// New creates a new instance of the encryptor service with the provided secret key.
//
// Parameters:
//   - secret: A byte slice representing the secret key used for encryption and decryption operations.
//
// Returns:
//   - Service: An implementation of the encryptor.Service interface, initialized with the provided secret key.
func New(secret []byte) Service {
	return &serviceImpl{
		ApiSecretAuthorization: secret,
	}
}

func (a *serviceImpl) Secret() []byte {
	return a.ApiSecretAuthorization
}

func (a *serviceImpl) Encrypt(value string) (string, error) {
	if value == "" {
		return "", fmt.Errorf("value cannot be empty")
	}
	return Encrypt(value, a.ApiSecretAuthorization)
}

func (a *serviceImpl) Decrypt(encrypted string) (string, error) {
	return Decrypt(encrypted, a.ApiSecretAuthorization)
}

func (a *serviceImpl) DecryptAll(encrypted ...string) ([]string, error) {
	var decrypted []string
	for _, value := range encrypted {
		decryptedValue, err := a.Decrypt(value)
		if err != nil {
			return nil, err
		}
		decrypted = append(decrypted, decryptedValue)
	}
	return decrypted, nil
}
