package security

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// GenerateApiSecretJWT creates a JWT token with the username and role
func (a *apiSecurityServiceImpl[T]) GenerateApiSecretJWT(jwtInfo ApiJWTInfo) (string, error) {
	secret := a.Secret()

	encryptedKey, err := a.Encrypt(jwtInfo.Key)
	if err != nil {
		return "", err
	}

	duration := time.Hour * jwtInfo.Expiration
	expiration := time.Now().Add(duration).Unix()
	claims := jwt.MapClaims{
		"client": jwtInfo.Client,
		"key":    encryptedKey,
		"exp":    expiration,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
