package impl

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/softwareplace/goserve/security/jwt/claims"
	"github.com/softwareplace/goserve/security/jwt/constants"
)

func NewClaims() claims.Claims {
	return &claimsImpl{}
}

type claimsImpl struct {
}

func (a *claimsImpl) Create(
	sub string,
	aud []string,
	exp int64,
	iat time.Time,
	iss string,
) jwt.MapClaims {
	claims := jwt.MapClaims{
		constants.SUB: sub,
		constants.AUD: aud,
		constants.EXP: exp,
		constants.IAT: iat.Unix(),
	}

	if iss != "" {
		claims[constants.ISS] = iss
	}

	return claims
}

func (a *claimsImpl) Get(token *jwt.Token) (jwt.MapClaims, bool) {
	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}
