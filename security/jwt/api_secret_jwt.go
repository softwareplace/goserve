package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func (a *serviceImpl[T]) GenerateApiSecretJWT(info Entry) (*Response, error) {
	secret := a.Secret()

	encryptedKey, err := a.Encrypt(info.Key)
	if err != nil {
		return nil, err
	}

	duration := info.Expiration
	expiration := time.Now().Add(duration).Unix()

	claims := jwt.MapClaims{
		"client": info.Client,
		"apiKey": encryptedKey,
		"exp":    expiration,
	}

	if info.Scopes != nil && len(info.Scopes) > 0 {
		var encryptedRoles []string
		for _, role := range info.Scopes {
			encryptedRole, err := a.Encrypt(role)
			if err != nil {
				return nil, err
			}
			encryptedRoles = append(encryptedRoles, encryptedRole)
		}

		claims["scope"] = encryptedRoles
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	return &Response{
		Token:   signedToken,
		Expires: int(expiration),
	}, err
}
