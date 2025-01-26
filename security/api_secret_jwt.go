package security

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func (a *apiSecurityServiceImpl[T]) GenerateApiSecretJWT(jwtInfo ApiJWTInfo) (*JwtResponse, error) {
	secret := a.Secret()

	encryptedKey, err := a.Encrypt(jwtInfo.Key)
	if err != nil {
		return nil, err
	}

	duration := jwtInfo.Expiration
	expiration := time.Now().Add(duration).Unix()

	claims := jwt.MapClaims{
		"client": jwtInfo.Client,
		"apiKey": encryptedKey,
		"exp":    expiration,
	}

	if jwtInfo.Scopes != nil && len(jwtInfo.Scopes) > 0 {
		var encryptedRoles []string
		for _, role := range jwtInfo.Scopes {
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
	return &JwtResponse{
		Token:   signedToken,
		Expires: int(expiration),
	}, err
}
