package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security/secret"
)

func TestLoginResourceHandlerValidation(t *testing.T) {
	t.Run("expects that return 401 when api secret is required for all resources but was not provided", func(t *testing.T) {
		testEnvSetup()
		defer testEnvCleanup()

		// Create a new request
		loginBody := strings.NewReader(`{"username": "my-username","password": "ynT9558iiMga&ayTVGs3Gc6ug1"}`)
		req, err := http.NewRequest("POST", "/login", loginBody)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		secretService := secret.New(
			secretProvider,
			securityService,
		)

		Default().
			LoginService(loginService).
			SecretService(secretService).
			SecurityService(securityService).
			NotFoundHandler().
			ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("expects that can login in when secret service was provided but public path skip was activated", func(t *testing.T) {
		testEnvSetup()
		defer testEnvCleanup()
		// Create a new request
		loginBody := strings.NewReader(`{"username": "my-username","password": "ynT9558iiMga&ayTVGs3Gc6ug1"}`)
		req, err := http.NewRequest("POST", "/login", loginBody)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		secretService := secret.New(
			secretProvider,
			securityService,
		).DisableForPublicPath(true)

		Default().
			LoginService(loginService).
			SecretService(secretService).
			SecurityService(securityService).
			NotFoundHandler().
			ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("expects that can get login response successfully when requires api secret and it was provided", func(t *testing.T) {
		testEnvSetup()
		defer testEnvCleanup()

		_ = os.Setenv("FULL_AUTHORIZATION", "true")

		loginBody := strings.NewReader(`{"username": "my-username","password": "ynT9558iiMga&ayTVGs3Gc6ug1"}`)
		req, err := http.NewRequest("POST", "/login", loginBody)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		key, err := getApiKey()
		if err != nil {
			t.Fatalf("Failed to get api key: %v", err)
		}

		req.Header.Set(goservectx.XApiKey, key)

		rr := httptest.NewRecorder()

		Default().
			LoginService(loginService).
			SecretService(secretService).
			SecurityService(securityService).
			NotFoundHandler().
			ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should response with bad request when user was not provided", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/login", nil)

		require.NoError(t, err)
		rr := httptest.NewRecorder()

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, "test")
		httpServer := create[*goservectx.DefaultContext]()

		httpServer.Login(ctx)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should response with forbidden when user is not authorized", func(t *testing.T) {
		testEnvSetup()
		defer testEnvCleanup()

		loginBody := strings.NewReader(`{"username": "my-username","password": "ynT9558iiMga&ayTVGs3Gc6ug1"}`)

		req, err := http.NewRequest("POST", "/login", loginBody)

		require.NoError(t, err)
		rr := httptest.NewRecorder()

		api := forTest(New[*goservectx.DefaultContext]().
			LoginService(loginService).
			SecurityService(securityService))

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, "test")

		api.Login(ctx)

		require.Equal(t, http.StatusOK, rr.Code)
	})
}
