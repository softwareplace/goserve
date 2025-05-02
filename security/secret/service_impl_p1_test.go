package secret

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/internal/service/login"
	"github.com/softwareplace/goserve/internal/service/provider"
	testutils "github.com/softwareplace/goserve/internal/utils"
	"github.com/softwareplace/goserve/security"
	model2 "github.com/softwareplace/goserve/security/jwt/model"
	"github.com/softwareplace/goserve/security/model"
	"github.com/softwareplace/goserve/security/router"
	"github.com/softwareplace/goserve/utils"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var secretFilePath = testutils.TestSecretFilePath()
var mockSecretKey = "DlJeR4%pPbB5Pr5cICMxg0xB"

type testSecurityService struct {
	security.Service[*goservectx.DefaultContext]
	testEncrypt func(value string) (string, error)
}

func (a *testSecurityService) Encrypt(value string) (string, error) {
	if a.testEncrypt != nil {
		return a.testEncrypt(value)
	}
	return a.Encrypt(value)
}

type testSecretHandlerImpl struct {
	apiSecretHandlerImpl[*goservectx.DefaultContext]
	tesApiSecretKeyValidation func(ctx *goservectx.Request[*goservectx.DefaultContext]) bool
}

func (a *testSecretHandlerImpl) ApiSecretKeyValidation(ctx *goservectx.Request[*goservectx.DefaultContext]) bool {
	if a.tesApiSecretKeyValidation != nil {
		return a.tesApiSecretKeyValidation(ctx)
	}
	return a.apiSecretHandlerImpl.ApiSecretKeyValidation(ctx)
}

func forTest(
	provider Provider[*goservectx.DefaultContext],
	service security.Service[*goservectx.DefaultContext],
	testEncrypt func(value string) (string, error),
) testSecretHandlerImpl {
	secretKey := utils.GetEnvOrDefault("API_PRIVATE_KEY", "")

	if secretKey == "" {
		log.Panicf("API_PRIVATE_KEY environment variable not set")
	}

	handler := apiSecretHandlerImpl[*goservectx.DefaultContext]{
		secretKey: secretKey,
		Service:   &testSecurityService{service, testEncrypt},
		Provider:  provider,
	}

	handler.InitAPISecretKey()

	return testSecretHandlerImpl{
		apiSecretHandlerImpl: handler,
	}
}

func init() {
	_ = os.Setenv("API_SECRET_KEY", mockSecretKey)
	_ = os.Setenv("JWT_CLAIMS_ENCRYPTION_ENABLED", "false")
	_ = os.Setenv("API_PRIVATE_KEY", secretFilePath)
}

func TestSecretImplValidation(t *testing.T) {

	t.Run("should run in panic when API_PRIVATE_KEY path was not set", func(t *testing.T) {
		_ = os.Unsetenv("API_PRIVATE_KEY")

		defer func() {
			_ = os.Setenv("API_PRIVATE_KEY", secretFilePath)
		}()

		var resultError error

		secretProvider := provider.NewSecretProvider()
		goserveerror.Handler(func() {
			_ = New[*goservectx.DefaultContext](
				secretProvider,
				security.New(login.NewPrincipalService()),
			)
		}, func(err error) {
			resultError = err
		})
		require.Error(t, resultError)
	})

	t.Run("should run in panic when API_PRIVATE_KEY path was provided but invalid", func(t *testing.T) {
		_ = os.Setenv("API_PRIVATE_KEY", "invalid-path-to-private-key.pem")

		defer func() {
			_ = os.Setenv("API_PRIVATE_KEY", secretFilePath)
		}()

		var resultError error

		secretProvider := provider.NewSecretProvider()
		goserveerror.Handler(func() {
			_ = New[*goservectx.DefaultContext](
				secretProvider,
				security.New(login.NewPrincipalService()),
			)
		}, func(err error) {
			resultError = err
		})
		require.Error(t, resultError)
	})

	t.Run("should return a new Service instance when API_PRIVATE_KEY exists", func(t *testing.T) {
		secretProvider := provider.NewSecretProvider()
		secretService := New[*goservectx.DefaultContext](
			secretProvider,
			security.New(login.NewPrincipalService()),
		)

		require.NotNil(t, secretService)
	})

	t.Run("should return expect secret", func(t *testing.T) {
		secretProvider := provider.NewSecretProvider()
		secretService := New[*goservectx.DefaultContext](
			secretProvider,
			security.New(login.NewPrincipalService()),
		)

		secret := secretService.Secret()

		require.Equal(t, string(secret), mockSecretKey)
	})

	t.Run("should call listener when api key successfully generated", func(t *testing.T) {
		secretProvider := provider.NewSecretProvider()
		secretService := New[*goservectx.DefaultContext](
			secretProvider,
			security.New(login.NewPrincipalService()),
		)

		req, err := http.NewRequest("POST", "/api-key/generate", nil)

		require.NoError(t, err)

		rr := httptest.NewRecorder()

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, goserveerror.HandlerWrapper)

		apiKeyEntryData := model.ApiKeyEntryData{
			ClientName: "test",
			ClientId:   "test",
			Expiration: 1000,
		}

		var expectedEntry *model.Entry

		secretProvider.TestOnGenerated = func(
			data model2.Response,
			jwtEntry model.Entry,
			ctx goservectx.SampleContext[*goservectx.DefaultContext],
		) {
			expectedEntry = &jwtEntry
		}

		secretService.Handler(ctx, apiKeyEntryData)

		require.NotNil(t, expectedEntry)
		require.NotNil(t, expectedEntry.Expiration)
		require.NotNil(t, expectedEntry.PublicKey)
		require.NotNil(t, expectedEntry.Roles)

		require.NotNil(t, expectedEntry.Key)
		require.Equal(t, expectedEntry.Key, "test")
	})

	t.Run("should return 200 when api key successfully generated ", func(t *testing.T) {
		secretProvider := provider.NewSecretProvider()
		secretService := New[*goservectx.DefaultContext](
			secretProvider,
			security.New(login.NewPrincipalService()),
		)

		req, err := http.NewRequest("POST", "/api-key/generate", nil)

		require.NoError(t, err)

		rr := httptest.NewRecorder()

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, goserveerror.HandlerWrapper)

		apiKeyEntryData := model.ApiKeyEntryData{
			ClientName: "test",
			ClientId:   "test",
			Expiration: 1000,
		}

		secretService.Handler(ctx, apiKeyEntryData)

		require.Equal(t, http.StatusOK, rr.Code)

		var response model2.Response

		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
	})

	t.Run("should return 500 when failed to GetJwtEntry ", func(t *testing.T) {
		secretProvider := provider.NewSecretProvider()
		secretService := New[*goservectx.DefaultContext](
			secretProvider,
			security.New(login.NewPrincipalService()),
		)

		req, err := http.NewRequest("POST", "/api-key/generate", nil)

		require.NoError(t, err)

		rr := httptest.NewRecorder()

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, goserveerror.HandlerWrapper)

		apiKeyEntryData := model.ApiKeyEntryData{
			ClientName: "test",
			ClientId:   "test",
			Expiration: 1000,
		}

		secretProvider.TestGetJwtEntry = func(
			apiKeyEntryData model.ApiKeyEntryData,
			_ *goservectx.Request[*goservectx.DefaultContext],
		) (model.Entry, error) {
			return model.Entry{}, fmt.Errorf("failed to GetJwtEntry")
		}

		secretService.Handler(ctx, apiKeyEntryData)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("should return 500 when public key is empty", func(t *testing.T) {
		secretProvider := provider.NewSecretProvider()
		secretService := New[*goservectx.DefaultContext](
			secretProvider,
			security.New(login.NewPrincipalService()),
		)

		req, err := http.NewRequest("POST", "/api-key/generate", nil)

		require.NoError(t, err)

		rr := httptest.NewRecorder()

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, goserveerror.HandlerWrapper)

		apiKeyEntryData := model.ApiKeyEntryData{
			ClientName: "test",
			ClientId:   "test",
			Expiration: 1000,
		}

		secretProvider.TestGetJwtEntry = func(
			apiKeyEntryData model.ApiKeyEntryData,
			_ *goservectx.Request[*goservectx.DefaultContext],
		) (model.Entry, error) {
			return model.Entry{
				Key:        "",
				Expiration: apiKeyEntryData.Expiration,
				Roles:      provider.MockScopes,
			}, nil
		}

		secretService.Handler(ctx, apiKeyEntryData)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("should return 500 when panic on generating public key", func(t *testing.T) {
		secretProvider := provider.NewSecretProvider()

		secretService := forTest(
			secretProvider,
			security.New(login.NewPrincipalService()),
			nil,
		)

		req, err := http.NewRequest("POST", "/api-key/generate", nil)

		require.NoError(t, err)

		rr := httptest.NewRecorder()

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, goserveerror.HandlerWrapper)

		apiKeyEntryData := model.ApiKeyEntryData{
			ClientName: "test",
			ClientId:   "test",
			Expiration: 1000,
		}

		secretService.secretKey = ""

		secretService.Handler(ctx, apiKeyEntryData)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("should return 500 when failed to generate public key", func(t *testing.T) {
		secretProvider := provider.NewSecretProvider()

		secretService := forTest(
			secretProvider,
			security.New(login.NewPrincipalService()),
			func(value string) (string, error) {
				return "", fmt.Errorf("failed to generate public key")
			},
		)

		req, err := http.NewRequest("POST", "/api-key/generate", nil)

		require.NoError(t, err)

		rr := httptest.NewRecorder()

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, goserveerror.HandlerWrapper)

		apiKeyEntryData := model.ApiKeyEntryData{
			ClientName: "test",
			ClientId:   "test",
			Expiration: 1000,
		}

		secretService.Handler(ctx, apiKeyEntryData)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("should return true when ignoreValidationForPublicPaths is enable and is a public path", func(t *testing.T) {
		secretProvider := provider.NewSecretProvider()

		secretService := New[*goservectx.DefaultContext](
			secretProvider,
			security.New(login.NewPrincipalService()),
		).DisableForPublicPath(true)

		req, err := http.NewRequest("POST", "http://localhost:8080/login", nil)
		rr := httptest.NewRecorder()

		router.AddOpenPath("POST::/login")
		require.NoError(t, err)

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, goserveerror.HandlerWrapper)

		require.Equal(t, true, secretService.HandlerSecretAccess(ctx))
	})

	t.Run("should return false when ignoreValidationForPublicPaths is true but is not a public path", func(t *testing.T) {
		secretProvider := provider.NewSecretProvider()

		secretService := New[*goservectx.DefaultContext](
			secretProvider,
			security.New(login.NewPrincipalService()),
		).DisableForPublicPath(true)

		req, err := http.NewRequest("POST", "http://localhost:8080/"+uuid.NewString(), nil)
		rr := httptest.NewRecorder()

		require.NoError(t, err)

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, goserveerror.HandlerWrapper)

		require.Equal(t, false, secretService.HandlerSecretAccess(ctx))
	})

	t.Run("should return true when ignoreValidationForPublicPaths is false and is not public path but is authorized", func(t *testing.T) {
		var expectedEntry model.Entry
		var data model2.Response

		generateJwtForTest(t, &expectedEntry, &data)

		secretProvider := provider.NewSecretProvider()
		secretService := New[*goservectx.DefaultContext](
			secretProvider,
			security.New(login.NewPrincipalService()),
		)

		req, err := http.NewRequest("POST", "/api-key/generate", nil)
		req.Header.Set(goservectx.XApiKey, data.JWT)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, goserveerror.HandlerWrapper)

		secretProvider.TestProviderGet = func(ctx *goservectx.Request[*goservectx.DefaultContext]) (string, error) {
			return *expectedEntry.PublicKey, nil
		}

		require.Equal(t, true, secretService.HandlerSecretAccess(ctx))
	})

}

func generateJwtForTest(t *testing.T, expectedEntry *model.Entry, targetData *model2.Response) {
	secretProvider := provider.NewSecretProvider()
	secretService := New[*goservectx.DefaultContext](
		secretProvider,
		security.New(login.NewPrincipalService()),
	)

	req, err := http.NewRequest("POST", "/api-key/generate", nil)

	require.NoError(t, err)

	rr := httptest.NewRecorder()

	ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, goserveerror.HandlerWrapper)

	apiKeyEntryData := model.ApiKeyEntryData{
		ClientName: "test",
		ClientId:   "test",
		Expiration: 1000,
	}

	secretProvider.TestOnGenerated = func(
		data model2.Response,
		jwtEntry model.Entry,
		ctx goservectx.SampleContext[*goservectx.DefaultContext],
	) {
		*targetData = data
		*expectedEntry = jwtEntry
	}

	secretService.Handler(ctx, apiKeyEntryData)

	require.NotNil(t, expectedEntry)
	require.NotNil(t, expectedEntry.Expiration)
	require.NotNil(t, expectedEntry.PublicKey)
	require.NotNil(t, expectedEntry.Roles)

	require.NotNil(t, expectedEntry.Key)
	require.Equal(t, expectedEntry.Key, "test")
}
