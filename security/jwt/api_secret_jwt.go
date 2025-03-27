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
	now := time.Now()
	expiration := now.Add(duration).Unix()

	claims := jwt.MapClaims{
		SUB: encryptedKey,
		EXP: expiration,
		IAT: now.Unix(),
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

		claims[AUD] = encryptedRoles
	}

	if issuer := a.Issuer(); issuer != "" {
		claims[ISS] = issuer
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	return &Response{
		JWT:      signedToken,
		Expires:  int(expiration),
		IssuedAt: int(now.Unix()),
	}, err
}
