package security

import "github.com/softwareplace/http-utils/utils"

func (a *apiSecurityServiceImpl[T]) Encrypt(value string) (string, error) {
	secret := a.Secret()
	return utils.Encrypt(value, secret)
}

func (a *apiSecurityServiceImpl[T]) Decrypt(encrypted string) (string, error) {
	secret := a.Secret()
	return utils.Decrypt(encrypted, secret)
}
