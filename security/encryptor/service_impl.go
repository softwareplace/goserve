package encryptor

import "github.com/softwareplace/http-utils/utils"

type ServiceImpl struct {
	secret []byte
}

func (a *ServiceImpl) Encrypt(value string) (string, error) {
	secret := a.secret
	return utils.Encrypt(value, secret)
}

func (a *ServiceImpl) Decrypt(encrypted string) (string, error) {
	return utils.Decrypt(encrypted, a.secret)
}
