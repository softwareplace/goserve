package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
)

type mockClaimsServiceImpl struct {
	returnStatus bool
	Claims
}

func newMockClaimsService(returnStatus bool) *mockClaimsServiceImpl {
	return &mockClaimsServiceImpl{
		Claims:       &claimsImpl{},
		returnStatus: returnStatus,
	}
}

func (m *mockClaimsServiceImpl) Get(token *jwt.Token) (jwt.MapClaims, bool) {
	return nil, m.returnStatus
}

func TestClaimsValidation(t *testing.T) {
	t.Run("should return true when is a valid jwt", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"
		service := New[*goservectx.DefaultContext](
			mockApiSecretKey,
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		response := generateJwt(t, 1*time.Minute)

		token, err := service.Parse(response.JWT)

		require.NoError(t, err)
		require.NotNil(t, token)

		claimsService := claimsImpl{}

		claims, isValid := claimsService.Get(token)

		require.True(t, true, isValid)
		require.NotEmpty(t, claims)
	})

	t.Run("should return expected jwt claims", func(t *testing.T) {
		claimsService := claimsImpl{}
		aud := []string{
			"test-role",
			"test-role2",
			"user-admin",
		}

		claims := claimsService.Create(
			"test",
			aud,
			time.Now().Add(1*time.Hour).Unix(),
			time.Now(),
			"test-issuer",
		)

		require.NotEmpty(t, claims)
		for _, value := range aud {
			require.True(t, contains(aud, value))
		}
	})
}
