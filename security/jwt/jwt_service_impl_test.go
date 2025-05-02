package jwt

import (
	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/security/encryptor"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type testPrincipalServiceImpl struct {
	status bool
}

func (d *testPrincipalServiceImpl) LoadPrincipal(ctx *goservectx.Request[*goservectx.DefaultContext]) bool {
	return d.status
}

func getDefaultCtx() *goservectx.DefaultContext {
	context := goservectx.NewDefaultCtx()
	context.SetRequesterId("gyo0V18QDj9Q1UWmZ2g7fc9sXrmlSthy3b8k9VO3MMv8dlEGtMtfIiPtJIUli0j")
	context.SetRoles("api:key:goserve-generator", "write:pets", "read:pets")
	return context
}

func TestJwtPrincipal(t *testing.T) {
	t.Run("should return true when principal loaded success", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		service := BaseService[*goservectx.DefaultContext]{
			Service: encryptor.New([]byte(mockApiSecretKey)),
			PService: &testPrincipalServiceImpl{
				status: true,
			},
			ErrorHandler: goserveerror.Default[*goservectx.DefaultContext](),
		}

		context := getDefaultCtx()

		jwtData, err := service.Generate(context, 15*time.Minute)

		require.NoError(t, err)
		require.NotEmpty(t, jwtData)

		req, err := http.NewRequest("POST", "/admin", nil)
		req.Header.Set(goservectx.XApiKey, jwtData.JWT)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, "test")

		require.True(t, service.Principal(ctx))
	})

	t.Run("should return false when principal loaded failed to load", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		service := BaseService[*goservectx.DefaultContext]{
			Service: encryptor.New([]byte(mockApiSecretKey)),
			PService: &testPrincipalServiceImpl{
				status: false,
			},
			ErrorHandler: goserveerror.Default[*goservectx.DefaultContext](),
		}

		context := getDefaultCtx()

		jwtData, err := service.Generate(context, 15*time.Minute)

		require.NoError(t, err)
		require.NotEmpty(t, jwtData)

		req, err := http.NewRequest("POST", "/admin", nil)
		req.Header.Set(goservectx.XApiKey, jwtData.JWT)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, "test")

		require.False(t, service.Principal(ctx))

		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestJwtValidation(t *testing.T) {
	_ = os.Setenv("JWT_ISSUER", "goserve-test-runner")

	t.Run("should return err when encryption failed", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN"

		service := New[*goservectx.DefaultContext](
			&testPrincipalServiceImpl{},
			mockApiSecretKey,
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		context := getDefaultCtx()

		jwtData, err := service.Generate(context, 15*time.Minute)

		require.Error(t, err)
		require.Nil(t, jwtData)
	})

	t.Run("should return err when sub is empty", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		service := New[*goservectx.DefaultContext](
			&testPrincipalServiceImpl{},
			mockApiSecretKey,
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		context := getDefaultCtx()
		context.SetRequesterId("")

		jwtData, err := service.Generate(context, 15*time.Minute)

		require.Error(t, err)
		require.Nil(t, jwtData)
	})

	t.Run("should return err when contains an empty role", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		service := New[*goservectx.DefaultContext](
			&testPrincipalServiceImpl{},
			mockApiSecretKey,
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		context := getDefaultCtx()
		context.SetRoles("")

		jwtData, err := service.Generate(context, 15*time.Minute)

		require.Error(t, err)
		require.Nil(t, jwtData)
	})

	t.Run("should generate a jwt with encrypted data", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		service := New[*goservectx.DefaultContext](
			&testPrincipalServiceImpl{},
			mockApiSecretKey,
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		context := getDefaultCtx()

		jwtData, err := service.Generate(context, 15*time.Minute)

		require.NoError(t, err)
		require.NotEmpty(t, jwtData)

		decode, err := service.Decode(jwtData.JWT)
		require.NoError(t, err)
		require.NotEmpty(t, decode)

		aud, ok := decode["aud"].([]interface{})

		require.True(t, ok)
		require.Equal(t, len(context.GetRoles()), len(aud))

		for _, value := range aud {
			require.False(t, contains(context.GetRoles(), value.(string)))
		}

		require.NotEqual(t, context.GetId(), decode["sub"])
		require.Equal(t, "goserve-test-runner", decode["iss"])

		require.NotEmpty(t, decode["iat"])
		require.NotEmpty(t, decode["exp"])
	})

	t.Run("should generate a jwt with original data when JWT_CLAIMS_ENCRYPTION_ENABLED provided as false", func(t *testing.T) {
		_ = os.Setenv("JWT_CLAIMS_ENCRYPTION_ENABLED", "false")

		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		service := New[*goservectx.DefaultContext](
			&testPrincipalServiceImpl{},
			mockApiSecretKey,
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		context := getDefaultCtx()

		jwtData, err := service.Generate(context, 15*time.Minute)

		require.NoError(t, err)
		require.NotEmpty(t, jwtData)

		decode, err := service.Decode(jwtData.JWT)
		require.NoError(t, err)
		require.NotEmpty(t, decode)

		aud, ok := decode["aud"].([]interface{})

		require.True(t, ok)
		require.Equal(t, len(context.GetRoles()), len(aud))

		for _, value := range aud {
			require.True(t, contains(context.GetRoles(), value.(string)))
		}

		require.NotEmpty(t, decode["sub"])
		require.Equal(t, context.GetId(), decode["sub"])
		require.Equal(t, "goserve-test-runner", decode["iss"])
		require.NotEmpty(t, decode["iat"])
		require.NotEmpty(t, decode["exp"])
	})

	t.Run("should return err given an invalid jwt", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		service := New[*goservectx.DefaultContext](
			&testPrincipalServiceImpl{},
			mockApiSecretKey,
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		decode, err := service.Decode("invalid")

		require.Error(t, err)
		require.Nil(t, decode)
	})

	t.Run("should return err given an valid jwt with invalid signature", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		service := New[*goservectx.DefaultContext](
			&testPrincipalServiceImpl{},
			mockApiSecretKey,
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		mockJwtInvalidSignature := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

		decode, err := service.Decode(mockJwtInvalidSignature)

		require.Error(t, err)
		require.Nil(t, decode)
	})
}

func contains[T comparable](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
