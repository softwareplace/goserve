package login

import (
	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security/encryptor"
)

// DefaultPasswordValidator is a generic type responsible for validating user passwords
// against their stored encrypted counterparts in principal contexts.
//
// T represents a type that implements the Principal interface.
// It ensures that the principal context contains methods to retrieve the encrypted password and other details.
type DefaultPasswordValidator[T goservectx.Principal] struct {
}

func (a *DefaultPasswordValidator[T]) IsValidPassword(loginEntryData User, principal T) bool {
	return encryptor.NewEncrypt(loginEntryData.Password).
		IsValidPassword(principal.EncryptedPassword())
}
