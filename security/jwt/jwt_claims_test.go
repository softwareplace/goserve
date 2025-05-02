package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type mockClaimsServiceImpl struct {
	returnStatus bool
}

func (m *mockClaimsServiceImpl) GetClaims(token *jwt.Token) (jwt.MapClaims, bool) {
	return nil, m.returnStatus
}
