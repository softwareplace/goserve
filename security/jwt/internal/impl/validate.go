package impl

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/softwareplace/goserve/security/jwt/validate"
)

type validateImpl struct {
	secret []byte
}

func (v *validateImpl) IsValid(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return v.secret, nil
	})

	return err == nil && token.Valid
}

func NewValidate(secret []byte) validate.Validate {
	return &validateImpl{secret: secret}
}
