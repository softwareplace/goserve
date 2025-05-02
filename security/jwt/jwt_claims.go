package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims interface {

	// GetClaims extracts and validates the JWT claims from the provided token.
	//
	// Parameters:
	//   - token: The JWT token from which claims are to be extracted.
	//
	// Returns:
	//   - jwt.MapClaims: A map containing the claims extracted from the token.
	//   - bool: A boolean value indicating whether the claims were successfully extracted and valid.
	GetClaims(token *jwt.Token) (jwt.MapClaims, bool)
}

func (a *claimsImpl[T]) GetClaims(token *jwt.Token) (jwt.MapClaims, bool) {
	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}
