package encryptor

import "github.com/softwareplace/http-utils/utils"

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
	return utils.Encrypt(value, a.ApiSecretAuthorization)
}

func (a *serviceImpl) Decrypt(encrypted string) (string, error) {
	return utils.Decrypt(encrypted, a.ApiSecretAuthorization)
}
