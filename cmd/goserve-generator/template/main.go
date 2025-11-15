package template

const GoServeMainTest = `package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/softwareplace/goserve/server"
	"github.com/stretchr/testify/require"

	"github.com/test-user/test-execution-validate-01/internal/adapter/handler"
	"github.com/test-user/test-execution-validate-01/internal/application"
)

func TestMockServer(t *testing.T) {
	t.Run("should return 400 status code access hello endpoint when required query not provided", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/api/${PROJECT}/v1/hello", nil)
		require.NoError(t, err, "❌ Failed to create request")

		rr := httptest.NewRecorder()

		server.New[*application.Principal]().
			ContextPath("/api/${PROJECT}/v1/").
			EmbeddedServer(handler.EmbeddedServer).
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("❌ handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		} else {
			log.Printf("✅ Expected status code %d", http.StatusBadRequest)
		}
	})

	t.Run("should return 200 status code access hello endpoint", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/api/${PROJECT}/v1/hello", nil)

		q := req.URL.Query()
		q.Add("username", "${PROJECT}")
		req.URL.RawQuery = q.Encode()

		require.NoError(t, err, "❌ Failed to create request")

		rr := httptest.NewRecorder()

		server.New[*application.Principal]().
			ContextPath("/api/${PROJECT}/v1/").
			EmbeddedServer(handler.EmbeddedServer).
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("❌ handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		} else {
			log.Printf("✅ Expected status code %d", http.StatusOK)
		}
	})

	t.Run("should return expected message by accessing hello endpoint", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/${PROJECT}/v1/hello", nil)

		q := req.URL.Query()
		q.Add("username", "Go Serve")
		req.URL.RawQuery = q.Encode()

		require.NoError(t, err, "❌ Failed to create request")

		rr := httptest.NewRecorder()

		server.New[*application.Principal]().
			ContextPath("/api/${PROJECT}/v1/").
			EmbeddedServer(handler.EmbeddedServer).
			ServeHTTP(rr, req)

		responseMessage := rr.Body.String()
		if strings.Contains(responseMessage, "Go Serve") {
			log.Printf("✅ Expected response body to contain %s", responseMessage)
		} else {
			t.Errorf("❌ Expected response body to contain 'Go Serve', but got: %s", responseMessage)
		}
	})
}
`

const GoServeMain = `package main

import (
	"github.com/softwareplace/goserve/logger"
	"github.com/softwareplace/goserve/server"

	"github.com/test-user/test-execution-validate-01/internal/adapter/handler"
	"github.com/test-user/test-execution-validate-01/internal/application"
	"github.com/test-user/test-execution-validate-01/internal/application/config"
)

func init() {
	// Setup log system. Using nested-logrus-formatter -> https://github.com/antonfisher/nested-logrus-formatter?tab=readme-ov-file
	// Reload log file target reference based on LOG_FILE_NAME_DATE_FORMAT
	logger.LogSetup()
}

func main() {
	server.New[*application.Principal]().
		Port(config.Port).
		ContextPath(config.ContextPath).
		SwaggerDocHandler(config.SwaggerFile).
		EmbeddedServer(handler.EmbeddedServer).
		StartServer()
}
`
