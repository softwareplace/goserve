package security

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func (a *apiSecurityServiceImpl[T]) GenerateApiSecretJWT(jwtInfo ApiJWTInfo) (string, error) {
	secret := a.Secret()

	encryptedKey, err := a.Encrypt(jwtInfo.Key)
	if err != nil {
		return "", err
	}

	duration := jwtInfo.Expiration
	expiration := time.Now().Add(duration).Unix()
	claims := jwt.MapClaims{
		"client": jwtInfo.Client,
		"apiKey": encryptedKey,
		"exp":    expiration,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
