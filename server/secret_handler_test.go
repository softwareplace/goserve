package server

import (
	goservectx "github.com/softwareplace/goserve/context"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApiSecretHandlerValidation(t *testing.T) {
	t.Run("should response with forbidden when user is not authorized", func(t *testing.T) {
		testEnvSetup()
		defer testEnvCleanup()

		body := strings.NewReader(`{"clientName": "test-client","clientId": "ynT9558iiMgaayTVGs3Gc6ug1"}`)

		req, err := http.NewRequest("POST", "/api-key/generate", body)

		require.NoError(t, err)
		rr := httptest.NewRecorder()

		api := forTest(New[*goservectx.DefaultContext]().
			SecretService(secretService))

		ctx := goservectx.Of[*goservectx.DefaultContext](rr, req, "test")

		api.ApiKeyGenerator(ctx)

		require.Equal(t, http.StatusOK, rr.Code)
	})
}
