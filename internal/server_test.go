package main

import (
	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/internal/service/apiservice"
	"github.com/softwareplace/goserve/internal/service/login"
	"github.com/softwareplace/goserve/internal/service/provider"
	"github.com/softwareplace/goserve/logger"
	"github.com/softwareplace/goserve/security"
	"github.com/softwareplace/goserve/security/secret"
	"github.com/softwareplace/goserve/server"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func init() {
	logger.LogReportCaller = true
	logger.LogSetup()
}

var (
	_ = os.Setenv("API_SECRET_KEY", "DlJeR4%pPbB5Pr5cICMxg0xB")
	_ = os.Setenv("JWT_CLAIMS_ENCRYPTION_ENABLED", "false")
	_ = os.Setenv("API_PRIVATE_KEY", "./secret/private.key")

	userPrincipalService = login.NewPrincipalService()
	securityService      = security.New(userPrincipalService)

	loginService   = login.NewLoginService(securityService)
	secretProvider = provider.NewSecretProvider()

	secretHandler = secret.New(
		secretProvider,
		securityService,
	)
)

func getApiKey() (string, error) {
	response, err := securityService.From(provider.MockJWTSub, provider.MockScopes, time.Minute*10)
	if err != nil {
		return "", err
	}
	return response.JWT, nil
}

func TestMockServer(t *testing.T) {
	t.Run("expects that can get login response successfully", func(t *testing.T) {
		// Create a new request
		loginBody := strings.NewReader(`{"username": "my-username","password": "ynT9558iiMga&ayTVGs3Gc6ug1"}`)
		req, err := http.NewRequest("POST", "/login", loginBody)

		//req.Header.Set("Content-Type", "application/json")

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		server.Default().
			LoginService(loginService).
			SecurityService(securityService).
			NotFoundHandler().
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("should return ok by request api health check", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		server.Default().
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("should return 404 by request api health check and it was disable", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		server.Default().
			HealthResourceEnabled(false).
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
	})

	t.Run("expects that return 401 when api secret is required for all resources but was not provided", func(t *testing.T) {
		// Create a new request
		loginBody := strings.NewReader(`{"username": "my-username","password": "ynT9558iiMga&ayTVGs3Gc6ug1"}`)
		req, err := http.NewRequest("POST", "/login", loginBody)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		secretProvider := provider.NewSecretProvider()
		secretService := secret.New(
			secretProvider,
			securityService,
		)

		server.Default().
			LoginService(loginService).
			SecretService(secretService).
			SecurityService(securityService).
			NotFoundHandler().
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusUnauthorized)
		}
	})

	t.Run("expects that can login in when secret service was provided but public path skip was activated", func(t *testing.T) {
		// Create a new request
		loginBody := strings.NewReader(`{"username": "my-username","password": "ynT9558iiMga&ayTVGs3Gc6ug1"}`)
		req, err := http.NewRequest("POST", "/login", loginBody)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		secretProvider := provider.NewSecretProvider()
		secretService := secret.New(
			secretProvider,
			securityService,
		).DisableForPublicPath(true)

		server.Default().
			LoginService(loginService).
			SecretService(secretService).
			SecurityService(securityService).
			NotFoundHandler().
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("expects that can get login response successfully when requires api secret and it was provided", func(t *testing.T) {
		// Create a new request

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

		server.Default().
			LoginService(loginService).
			SecretService(secretHandler).
			SecurityService(securityService).
			NotFoundHandler().
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("expects that return default not found when a custom was not provided", func(t *testing.T) {
		// Create a new request
		req, err := http.NewRequest("POST", "/not-found", nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		server.Default().
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}

		if strings.Contains(rr.Body.String(), "404 page not found") {
			t.Log("Response body contains '404 page not found'")
		} else {
			t.Errorf("Expected response body to contain '404 page not found', but got: %s", rr.Body.String())
		}

	})

	t.Run("expects that return custom not found when a custom was provided", func(t *testing.T) {
		// Create a new request
		req, err := http.NewRequest("POST", "/not-found", nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		server.Default().
			CustomNotFoundHandler(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte("Custom 404 Page"))
			}).
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}

		if strings.Contains(rr.Body.String(), "Custom 404 Page") {
			t.Log("Response body contains 'Custom 404 Page'")
		} else {
			t.Errorf("Expected response body to contain 'Custom 404 Page', but got: %s", rr.Body.String())
		}
	})

	t.Run("expects that return swagger resource when swagger was defined and using the default not found handler", func(t *testing.T) {
		// Create a new request
		req, err := http.NewRequest("GET", "/", nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		server.Default().
			EmbeddedServer(apiservice.Register).
			SwaggerDocHandler("./resource/pet-store.yaml").
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMovedPermanently {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusMovedPermanently)
		}

		if strings.Contains(rr.Body.String(), "<a href=\"/swagger/index.html\">Moved Permanently</a>.") {
			t.Log("Response body contains '<a href=\"/swagger/index.html\">Moved Permanently</a>.'")
		} else {
			t.Errorf("Expected response body to contain '<a href=\"/swagger/index.html\">Moved Permanently</a>.', but got: %s", rr.Body.String())
		}
	})
}
