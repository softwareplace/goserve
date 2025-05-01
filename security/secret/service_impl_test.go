package secret

import (
	"encoding/json"
	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/internal/service/login"
	"github.com/softwareplace/goserve/internal/service/provider"
	"github.com/softwareplace/goserve/internal/utils"
	"github.com/softwareplace/goserve/security"
	"github.com/softwareplace/goserve/security/jwt"
	"github.com/softwareplace/goserve/security/model"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var secretFilePath = utils.TestSecretFilePath()
var mockSecretKey = "DlJeR4%pPbB5Pr5cICMxg0xB"

func init() {
	_ = os.Setenv("API_SECRET_KEY", mockSecretKey)
	_ = os.Setenv("JWT_CLAIMS_ENCRYPTION_ENABLED", "false")
	_ = os.Setenv("API_PRIVATE_KEY", secretFilePath)
}

func TestSecretImplValidation(t *testing.T) {
	secretProvider := provider.NewSecretProvider()

	t.Run("should run in panic when API_PRIVATE_KEY path was not set", func(t *testing.T) {
		_ = os.Unsetenv("API_PRIVATE_KEY")

		defer func() {
			_ = os.Setenv("API_PRIVATE_KEY", secretFilePath)
		}()

		var resultError error

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
		secretService := New[*goservectx.DefaultContext](
			secretProvider,
			security.New(login.NewPrincipalService()),
		)

		require.NotNil(t, secretService)
	})

	t.Run("should return expect secret", func(t *testing.T) {
		secretService := New[*goservectx.DefaultContext](
			secretProvider,
			security.New(login.NewPrincipalService()),
		)

		secret := secretService.Secret()

		require.Equal(t, string(secret), mockSecretKey)
	})

	t.Run("should call listener when api key successfully generated", func(t *testing.T) {
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

		defer func() {
			secretProvider.Callback = nil
		}()

		callBack := utils.TestCallBack[model.Entry](func(entry model.Entry) {
			expectedEntry = &entry
		})

		secretProvider.Callback = &callBack
		secretService.Handler(ctx, apiKeyEntryData)

		require.NotNil(t, expectedEntry)
		require.NotNil(t, expectedEntry.Expiration)
		require.NotNil(t, expectedEntry.PublicKey)
		require.NotNil(t, expectedEntry.Roles)

		require.NotNil(t, expectedEntry.Key)
		require.Equal(t, expectedEntry.Key, "test")
	})

	t.Run("should return 200 when api key successfully generated ", func(t *testing.T) {
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

		require.Equal(t, rr.Code, http.StatusOK)

		var response jwt.Response

		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
	})
}
