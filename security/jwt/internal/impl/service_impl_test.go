package impl

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/internal/service/testencryptor"
	"github.com/softwareplace/goserve/security/encryptor"
	"github.com/softwareplace/goserve/security/jwt/claims"
	"github.com/softwareplace/goserve/security/jwt/constants"
	"github.com/softwareplace/goserve/security/jwt/response"
)

type mockClaimsServiceImpl struct {
	returnStatus bool
	claims.Claims
}

func newMockClaimsService(returnStatus bool) *mockClaimsServiceImpl {
	return &mockClaimsServiceImpl{
		Claims:       NewClaims(),
		returnStatus: returnStatus,
	}
}

func (m *mockClaimsServiceImpl) Get(token *jwt.Token) (jwt.MapClaims, bool) {
	return nil, m.returnStatus
}

func TestErrorHandlerValidation(t *testing.T) {
	t.Run("should return false when claims successful extracted", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl[*goservectx.DefaultContext](
			NewClaims(),
			NewValidate(secret),
			encryptor.New(secret),
			nil,
		)

		req, err := http.NewRequest("POST", "/admin", nil)

		require.NoError(t, err)
		rr := httptest.NewRecorder()
		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, "test")

		service.HandlerErrorOrElse(ctx, nil, goserveerror.ExtractClaimsError, func() {
			ctx.InternalServerError("test")
		})

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestExtractJWTClaims(t *testing.T) {
	t.Run("should return false when claims successful extracted", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl[*goservectx.DefaultContext](
			NewClaims(),
			NewValidate(secret),
			encryptor.New(secret),
			nil,
		)

		context := getDefaultCtx()

		jwtData, err := service.Generate(context, 15*time.Minute)

		require.NoError(t, err)
		require.NotEmpty(t, jwtData)

		req, err := http.NewRequest("POST", "/admin", nil)
		req.Header.Set(goservectx.XApiKey, jwtData.JWT)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, "test")

		isSuccess := service.ExtractJWTClaims(ctx)
		require.False(t, isSuccess)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("should return false when got a expired jwt", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl(
			NewClaims(),
			NewValidate(secret),
			encryptor.New(secret),
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		context := getDefaultCtx()

		jwtData, err := service.Generate(context, 1*time.Second)

		require.NoError(t, err)
		require.NotEmpty(t, jwtData)

		time.Sleep(2 * time.Second)

		req, err := http.NewRequest("POST", "/admin", nil)
		req.Header.Set(goservectx.Authorization, jwtData.JWT)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, "test")

		require.False(t, service.ExtractJWTClaims(ctx))
	})

	t.Run("should return false when failed to decrypt claims", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		testEncryptor := testencryptor.New(encryptor.New([]byte(mockApiSecretKey)))

		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl(
			NewClaims(),
			NewValidate(secret),
			testEncryptor,
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		context := getDefaultCtx()

		jwtData, err := service.Generate(context, 1*time.Minute)

		require.NoError(t, err)
		require.NotEmpty(t, jwtData)

		req, err := http.NewRequest("POST", "/admin", nil)
		req.Header.Set(goservectx.Authorization, jwtData.JWT)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, "test")

		testEncryptor.TestDecrypt = func(encrypted string) (string, error) {
			return "", errors.New("failed to decrypt")
		}

		isSuccessful := service.ExtractJWTClaims(ctx)
		require.False(t, isSuccessful)
	})

	t.Run("should return false when failed to get claims", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		testEncryptor := testencryptor.New(encryptor.New([]byte(mockApiSecretKey)))

		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl(
			newMockClaimsService(false),
			NewValidate(secret),
			testEncryptor,
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		context := getDefaultCtx()

		jwtData, err := service.Generate(context, 1*time.Minute)

		require.NoError(t, err)
		require.NotEmpty(t, jwtData)

		req, err := http.NewRequest("POST", "/admin", nil)
		req.Header.Set(goservectx.Authorization, jwtData.JWT)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, "test")

		isSuccessful := service.ExtractJWTClaims(ctx)
		require.False(t, isSuccessful)
	})

	t.Run("should return true when claims extraction was successful", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl(
			NewClaims(),
			NewValidate(secret),
			encryptor.New(secret),
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		context := getDefaultCtx()

		jwtData := generateJwt(t, 15*time.Minute)

		req, err := http.NewRequest("POST", "/admin", nil)
		req.Header.Set(goservectx.Authorization, jwtData.JWT)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, "test")

		isSuccess := service.ExtractJWTClaims(ctx)
		require.True(t, isSuccess)
		require.Equal(t, context.GetId(), ctx.AccessId)
	})
}

func TestDecode(t *testing.T) {
	t.Run("should return nil when given an jwt with invalid signature", func(t *testing.T) {

		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		testEncryptor := testencryptor.New(encryptor.New([]byte("iBID8F32zkN1a0d4hCdm4gVg")))

		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl(
			NewClaims(),
			NewValidate(secret),
			testEncryptor,
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		jwtData := generateJwt(t, 15*time.Minute)
		decode, err := service.Decode(jwtData.JWT)
		require.Nil(t, decode)
		require.Error(t, err)
	})
}

func TestJwtValidation(t *testing.T) {
	_ = os.Setenv("JWT_ISSUER", "goserve-test-runner")

	t.Run("should return err when encryption failed", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN"

		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl(
			NewClaims(),
			NewValidate(secret),
			encryptor.New(secret),
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		context := getDefaultCtx()

		jwtData, err := service.Generate(context, 15*time.Minute)

		require.Error(t, err)
		require.Nil(t, jwtData)
	})

	t.Run("should return err when sub is empty", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl(
			NewClaims(),
			NewValidate(secret),
			encryptor.New(secret),
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

		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl(
			NewClaims(),
			NewValidate(secret),
			encryptor.New(secret),
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

		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl(
			NewClaims(),
			NewValidate(secret),
			encryptor.New(secret),
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		context := getDefaultCtx()

		jwtData, err := service.Generate(context, 15*time.Minute)

		require.NoError(t, err)
		require.NotEmpty(t, jwtData)

		decode, err := service.Decode(jwtData.JWT)
		require.NoError(t, err)
		require.NotEmpty(t, decode)

		aud, ok := decode[constants.AUD].([]interface{})

		require.True(t, ok)
		require.Equal(t, len(context.GetRoles()), len(aud))

		for _, value := range aud {
			require.False(t, contains(context.GetRoles(), value.(string)))
		}

		require.NotEqual(t, context.GetId(), decode[constants.SUB])
		require.Equal(t, "goserve-test-runner", decode[constants.ISS])

		require.NotEmpty(t, decode[constants.IAT])
		require.NotEmpty(t, decode[constants.EXP])
	})

	t.Run("should return a jwt decrypted claims", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl(
			NewClaims(),
			NewValidate(secret),
			encryptor.New(secret),
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		context := getDefaultCtx()

		jwtData, err := service.Generate(context, 15*time.Minute)

		require.NoError(t, err)
		require.NotEmpty(t, jwtData)

		decode, err := service.Decrypted(jwtData.JWT)
		require.NoError(t, err)
		require.NotEmpty(t, decode)

		aud, ok := decode[constants.AUD].([]string)

		require.True(t, ok)
		require.Equal(t, len(context.GetRoles()), len(aud))

		for _, value := range aud {
			require.True(t, contains(context.GetRoles(), value))
		}

		require.Equal(t, context.GetId(), decode[constants.SUB])
		require.Equal(t, "goserve-test-runner", decode[constants.ISS])

		require.NotEmpty(t, decode[constants.IAT])
		require.NotEmpty(t, decode[constants.EXP])
	})

	t.Run("should generate a jwt with original data when JWT_CLAIMS_ENCRYPTION_ENABLED provided as false", func(t *testing.T) {
		_ = os.Setenv("JWT_CLAIMS_ENCRYPTION_ENABLED", "false")

		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl(
			NewClaims(),
			NewValidate(secret),
			encryptor.New(secret),
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		context := getDefaultCtx()

		jwtData, err := service.Generate(context, 15*time.Minute)

		require.NoError(t, err)
		require.NotEmpty(t, jwtData)

		decode, err := service.Decode(jwtData.JWT)
		require.NoError(t, err)
		require.NotEmpty(t, decode)

		aud, ok := decode[constants.AUD].([]interface{})

		require.True(t, ok)
		require.Equal(t, len(context.GetRoles()), len(aud))

		for _, value := range aud {
			require.True(t, contains(context.GetRoles(), value.(string)))
		}

		require.NotEmpty(t, decode[constants.SUB])
		require.Equal(t, context.GetId(), decode[constants.SUB])
		require.Equal(t, "goserve-test-runner", decode[constants.ISS])
		require.NotEmpty(t, decode[constants.IAT])
		require.NotEmpty(t, decode[constants.EXP])
	})

	t.Run("should return err given an invalid jwt", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl(
			NewClaims(),
			NewValidate(secret),
			encryptor.New(secret),
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		decode, err := service.Decode("invalid")

		require.Error(t, err)
		require.Nil(t, decode)
	})

	t.Run("should return err given an valid jwt with invalid signature", func(t *testing.T) {
		mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"

		secret := []byte(mockApiSecretKey)

		service := NewJwtServiceImpl(
			NewClaims(),
			NewValidate(secret),
			encryptor.New(secret),
			goserveerror.Default[*goservectx.DefaultContext](),
		)

		mockJwtInvalidSignature := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

		decode, err := service.Decode(mockJwtInvalidSignature)

		require.Error(t, err)
		require.Nil(t, decode)
	})
}

func generateJwt(t *testing.T, duration time.Duration) *response.Response {
	mockApiSecretKey := "iBID8F32zkN1a0d4hCdm4gVS"
	secret := []byte(mockApiSecretKey)

	service := NewJwtServiceImpl(
		NewClaims(),
		NewValidate(secret),
		encryptor.New(secret),
		goserveerror.Default[*goservectx.DefaultContext](),
	)

	context := getDefaultCtx()

	jwtData, err := service.Generate(context, duration)
	require.NoError(t, err)
	require.NotEmpty(t, jwtData)
	return jwtData
}

func contains[T comparable](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
