package server

import (
	"github.com/softwareplace/goserve/logger"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	logger.LogReportCaller = true
	logger.LogSetup()
}

func TestMockServer(t *testing.T) {
	t.Run("expects that return default not found when a custom was not provided", func(t *testing.T) {
		// Create a new request
		req, err := http.NewRequest("POST", "/not-found", nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		Default().
			ServeHTTP(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code)

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

		Default().
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

}
