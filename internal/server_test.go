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
)

func init() {
	logger.LogReportCaller = true
	logger.LogSetup()
}

var (
	_ = os.Setenv("API_SECRET_KEY", "ue1pUOtCGaYS7Z1DLJ80nFtZ")

	userPrincipalService = login.NewPrincipalService()
	securityService      = security.New(userPrincipalService)

	loginService   = login.NewLoginService(securityService)
	secretProvider = provider.NewSecretProvider()

	secretHandler = secret.New(
		"./secret/private.key",
		secretProvider,
		securityService,
	)

	apiSecret = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiRWZUZElvbm8vd3o2ZUxOWG5PajZOQ2lOUnpqTWNMSklYdWlEeUt5d1ZEST0iLCI1elA0cWs1a2Q2aUtCekZ2Y2NiTUI3OW5EbUEwczgrY0dsMHVZT2s4MUE5cCIsIjdMVnZDVTlXbVl2SVY2OU1sTHdIZHpXb0hlV0VSSlBpQ1E9PSIsInNNem8vYjlUTGVHMVBwUjFkYkV5MGhmRC9vbHZkalZpeVIwPSIsIngzODhLdTkxdUJHTncwckp1MHcyRVhIR0JZajVKVUVaZFBuV2g0b1JyMk1rIl0sImV4cCI6NTE5OTA5ODAxMywiaWF0IjoxNzQzMDk4MDEzLCJpc3MiOiJnb3NlcnZlci1leGFtcGxlIiwic3ViIjoibFdQKzdHTjNzZjhoNVZXcVRyaTBUM0RaSHNaYmEvWWcwenV4TWhKK0o4Mkw2R0FHelRkUFl6N2hGV0doWkhBYiJ9.6-Z4W5np8uXLuQJttd9BOvuG7iG9EFC8RsTL2fB0OqU"
)

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
			"./secret/private.key",
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
			"./secret/private.key",
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
		loginBody := strings.NewReader(`{"username": "my-username","password": "ynT9558iiMga&ayTVGs3Gc6ug1"}`)
		req, err := http.NewRequest("POST", "/login", loginBody)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set(goservectx.XApiKey, apiSecret)

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
