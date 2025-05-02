package testencryptor

import "github.com/softwareplace/goserve/security/encryptor"

type Service struct {
	Original    encryptor.Service
	TesSecret   func() []byte
	TestEncrypt func(value string) (string, error)
	TestDecrypt func(encrypted string) (string, error)
}

func (t Service) Secret() []byte {
	if t.TesSecret != nil {
		return t.TesSecret()
	}
	return t.Original.Secret()
}

func (t Service) Encrypt(value string) (string, error) {
	if t.TestEncrypt != nil {
		return t.TestEncrypt(value)
	}
	return t.Original.Encrypt(value)
}

func (t Service) Decrypt(encrypted string) (string, error) {
	if t.TestDecrypt != nil {
		return t.TestDecrypt(encrypted)
	}
	return t.Original.Decrypt(encrypted)
}

func New(service encryptor.Service) *Service {
	return &Service{
		Original: service,
	}
}
