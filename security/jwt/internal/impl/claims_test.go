package impl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/security/encryptor"
)

func TestClaimsValidation(t *testing.T) {
	t.Run("should return true when is a valid jwt", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"
		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl(
			NewClaims(),
			NewValidate(secret),
			encryptor.New(secret),
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		response := generateJwt(t, 1*time.Minute)

		token, err := service.Parse(response.JWT)

		require.NoError(t, err)
		require.NotNil(t, token)

		claimsService := NewClaims()

		claims, isValid := claimsService.Get(token)

		require.True(t, true, isValid)
		require.NotEmpty(t, claims)
	})

	t.Run("should return expected jwt claims", func(t *testing.T) {
		claimsService := NewClaims()
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
