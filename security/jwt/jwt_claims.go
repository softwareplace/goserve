package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims interface {

	// Get extracts and validates the JWT claims from the provided token.
	//
	// Parameters:
	//   - token: The JWT token from which claims are to be extracted.
	//
	// Returns:
	//   - jwt.MapClaims: A map containing the claims extracted from the token.
	//   - bool: A boolean value indicating whether the claims were successfully extracted and valid.
	Get(token *jwt.Token) (jwt.MapClaims, bool)

	// Create generates a set of JWT claims based on the provided inputs.
	//
	// Parameters:
	//   - sub: The subject of the JWT (usually the user or entity making the request).
	//   - aud: A list of roles or audiences associated with the JWT.
	//   - exp: The expiration time of the JWT, represented as a Unix timestamp.
	//   - iat: The issued-at time of the JWT, represented as a time.Time object.
	//   - iss: The issuer of the JWT. If not empty, it will be included in the claims.
	//
	// Returns:
	//   - jwt.MapClaims: A map containing the generated claims.
	Create(
		sub string,
		aud []string,
		exp int64,
		iat time.Time,
		iss string,
	) jwt.MapClaims
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
		SUB: sub,
		AUD: aud,
		EXP: exp,
		IAT: iat.Unix(),
	}

	if iss != "" {
		claims[ISS] = iss
	}

	return claims
}

func (a *claimsImpl) Get(token *jwt.Token) (jwt.MapClaims, bool) {
	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}
