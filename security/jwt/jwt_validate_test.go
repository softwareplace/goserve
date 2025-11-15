package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
)

func TestIsValid(t *testing.T) {
	t.Run("should return false when token already expired", func(t *testing.T) {
		jwtData := generateJwt(t, 1*time.Second)

		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		service := New[*goservectx.DefaultContext](
			mockApiSecretKey,
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		time.Sleep(2 * time.Second)

		isSuccess := service.IsValid(jwtData.JWT)

		require.False(t, isSuccess)
	})

	t.Run("should return false when token has no a valid signature", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"
		service := New[*goservectx.DefaultContext](
			mockApiSecretKey,
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		isSuccess := service.IsValid("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")

		require.False(t, isSuccess)
	})

	t.Run("should return true when token is valid", func(t *testing.T) {
		jwtData := generateJwt(t, 1*time.Minute)

		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		service := New[*goservectx.DefaultContext](
			mockApiSecretKey,
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		isSuccess := service.IsValid(jwtData.JWT)

		require.True(t, isSuccess)
	})
}
