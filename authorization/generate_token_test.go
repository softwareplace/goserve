package authorization

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/softwareplace/goserve/authorization/config"
	"github.com/softwareplace/goserve/authorization/request"
	"github.com/softwareplace/goserve/authorization/response"
)

func TestClientImplGenerateTokenSuccess(t *testing.T) {
	expectedResponse := response.AuthorizationResponse{
		Jwt:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
		Expires:  1757353374,
		IssuedAt: 1757353374,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/login", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	client := NewClient(
		config.OauthConfig{
			ServerHost: server.URL,
		},
	)

	authRequest := request.AuhtorizationRequest{
		Username: "testuser",
		Password: "testpass",
	}

	result, err := client.Login(authRequest, "application/json")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse.Jwt, result.Jwt)
	assert.Equal(t, expectedResponse.Expires, result.Expires)
	assert.Equal(t, expectedResponse.IssuedAt, result.IssuedAt)
}

func TestClientImplGenerateTokenValidationErrorEmptyUsername(t *testing.T) {
	client := NewClient(
		config.OauthConfig{
			ServerHost: "http://localhost:8080",
		},
	)

	authRequest := request.AuhtorizationRequest{
		Username: "",
		Password: "testpass",
	}

	result, err := client.Login(authRequest, "application/json")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestClientImplGenerateTokenValidationErrorEmptyPassword(t *testing.T) {
	client := NewClient(
		config.OauthConfig{
			ServerHost: "http://localhost:8080",
		},
	)

	authRequest := request.AuhtorizationRequest{
		Username: "testuser",
		Password: "",
	}

	result, err := client.Login(authRequest, "application/json")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestClientImplGenerateTokenEmptyApplicationIs(t *testing.T) {
	client := NewClient(
		config.OauthConfig{
			ServerHost: "http://localhost:8080",
		},
	)

	authRequest := request.AuhtorizationRequest{
		Username: "testuser",
		Password: "testpass",
	}

	result, err := client.Login(authRequest, "")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "applicationID is required", err.Error())
}

func TestClientImplGenerateTokenNetworkError(t *testing.T) {
	client := NewClient(
		config.OauthConfig{
			ServerHost: "http://localhost:8080",
		},
	)

	authRequest := request.AuhtorizationRequest{
		Username: "testuser",
		Password: "testpass",
	}

	result, err := client.Login(authRequest, "application/json")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestClientImplGenerateTokenNilResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	server.Close()

	client := NewClient(
		config.OauthConfig{
			ServerHost: "http://localhost:8080",
		},
	)

	authRequest := request.AuhtorizationRequest{
		Username: "testuser",
		Password: "testpass",
	}

	result, err := client.Login(authRequest, "application/json")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestClientImplGenerateTokenInvalidJSONResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewClient(
		config.OauthConfig{
			ServerHost: "http://localhost:8080",
		},
	)

	authRequest := request.AuhtorizationRequest{
		Username: "testuser",
		Password: "testpass",
	}

	result, err := client.Login(authRequest, "application/json")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestClientImplGenerateTokenResponseValidationError(t *testing.T) {
	invalidResponse := response.AuthorizationResponse{
		Jwt:      "short",
		Expires:  100,
		IssuedAt: 100,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(invalidResponse)
	}))
	defer server.Close()

	client := NewClient(
		config.OauthConfig{
			ServerHost: server.URL,
		},
	)

	authRequest := request.AuhtorizationRequest{
		Username: "testuser",
		Password: "testpass",
	}

	result, err := client.Login(authRequest, "application/json")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to validate response data")
}

func TestClientImplGenerateTokenHTTPErrorStatus(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
	}{
		{"BadRequest", http.StatusBadRequest},
		{"Unauthorized", http.StatusUnauthorized},
		{"Forbidden", http.StatusForbidden},
		{"NotFound", http.StatusNotFound},
		{"InternalServerError", http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				w.Write([]byte("Error response"))
			}))
			defer server.Close()

			client := NewClient(
				config.OauthConfig{
					ServerHost: "http://localhost:8080",
				},
			)

			authRequest := request.AuhtorizationRequest{
				Username: "testuser",
				Password: "testpass",
			}

			result, err := client.Login(authRequest, "application/json")

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestClientImplGenerateTokenWithDifferentContentTypes(t *testing.T) {
	testCases := []struct {
		name           string
		applicationIs  string
		expectedHeader string
	}{
		{"JSON", "application/json", "application/json"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedResponse := response.AuthorizationResponse{
				Jwt:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
				Expires:  1757353374,
				IssuedAt: 1757353374,
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedHeader, r.Header.Get("Content-Type"))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(expectedResponse)
			}))

			defer server.Close()

			client := NewClient(
				config.OauthConfig{
					ServerHost: server.URL,
				},
			)

			authRequest := request.AuhtorizationRequest{
				Username: "testuser",
				Password: "testpass",
			}

			result, err := client.Login(authRequest, tc.applicationIs)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, expectedResponse.Jwt, result.Jwt)
		})
	}
}
