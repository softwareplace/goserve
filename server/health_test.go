package server

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthResourceHandlerTest(t *testing.T) {
	t.Run("expects that can get login response successfully", func(t *testing.T) {
		testEnvSetup()
		defer testEnvCleanup()
		// Create a new request
		loginBody := strings.NewReader(`{"username": "my-username","password": "ynT9558iiMga&ayTVGs3Gc6ug1"}`)
		req, err := http.NewRequest("POST", "/login", loginBody)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		Default().
			LoginService(loginService).
			SecurityService(securityService).
			NotFoundHandler().
			ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return ok by request api health check", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		Default().
			ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return 404 by request api health check and it was disable", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		Default().
			HealthResourceEnabled(false).
			ServeHTTP(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code)
	})
}
