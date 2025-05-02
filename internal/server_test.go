package main

import (
	"github.com/softwareplace/goserve/internal/service/apiservice"
	"github.com/softwareplace/goserve/logger"
	"github.com/softwareplace/goserve/server"
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
