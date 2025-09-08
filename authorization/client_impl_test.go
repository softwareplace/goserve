package authorization

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/softwareplace/goserve/request"
	"github.com/stretchr/testify/assert"
)

func TestClientImplCheckTokenSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &clientImpl{
		oauthConfig: OauthConfig{
			ServerHost: server.URL,
		},
	}

	input := Input{
		Headers: map[string]string{
			"Authorization": "Bearer token123",
			"Content-Type":  "application/json",
		},
		QueryParams: map[string][]string{
			"scope": {"read"},
			"type":  {"access_token"},
		},
	}

	result, err := client.CheckToken(input)

	assert.NoError(t, err)
	assert.True(t, result)
}

func TestClientImplCheckTokenFailureNonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := &clientImpl{
		oauthConfig: OauthConfig{
			ServerHost: server.URL,
		},
	}

	input := Input{
		Headers: map[string]string{
			"Authorization": "Bearer invalid_token",
		},
		QueryParams: map[string][]string{},
	}

	result, err := client.CheckToken(input)

	assert.NoError(t, err)
	assert.False(t, result)
}

func TestClientImplCheckTokenWithEmptyInput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &clientImpl{
		oauthConfig: OauthConfig{
			ServerHost: server.URL,
		},
	}

	input := Input{
		Headers:     map[string]string{},
		QueryParams: map[string][]string{},
	}

	result, err := client.CheckToken(input)

	assert.NoError(t, err)
	assert.True(t, result)
}

func TestClientImplChecktokenCustomSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &clientImpl{
		oauthConfig: OauthConfig{
			ServerHost: "http://example.com",
		},
	}

	config := &request.Config{
		Host:    server.URL,
		Path:    "/validate",
		Headers: map[string]string{"Authorization": "Bearer token123"},
		Query:   map[string][]string{"scope": {"read"}},
	}

	result, err := client.ChecktokenCustom(config)

	assert.NoError(t, err)
	assert.True(t, result)
}

func TestClientImplChecktokenCustomFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := &clientImpl{
		oauthConfig: OauthConfig{
			ServerHost: "http://example.com",
		},
	}

	config := &request.Config{
		Host:    server.URL,
		Path:    "/validate",
		Headers: map[string]string{"Authorization": "Bearer invalid_token"},
		Query:   map[string][]string{},
	}

	result, err := client.ChecktokenCustom(config)

	assert.NoError(t, err)
	assert.False(t, result)
}

func TestClientImplValidateNetworkError(t *testing.T) {
	client := &clientImpl{
		oauthConfig: OauthConfig{
			ServerHost: "http://invalid-host-that-does-not-exist.com",
		},
	}

	config := &request.Config{
		Host:    "http://invalid-host-that-does-not-exist.com",
		Path:    "/validate",
		Headers: map[string]string{},
		Query:   map[string][]string{},
	}

	result, err := client.validate(config)

	assert.Error(t, err)
	assert.False(t, result)
}

func TestClientImplValidateStatusOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &clientImpl{
		oauthConfig: OauthConfig{
			ServerHost: server.URL,
		},
	}

	config := &request.Config{
		Host:    server.URL,
		Path:    "/validate",
		Headers: map[string]string{},
		Query:   map[string][]string{},
	}

	result, err := client.validate(config)

	assert.NoError(t, err)
	assert.True(t, result)
}

func TestClientImplValidateStatusNotOK(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
	}{
		{"Unauthorized", http.StatusUnauthorized},
		{"Forbidden", http.StatusForbidden},
		{"NotFound", http.StatusNotFound},
		{"InternalServerError", http.StatusInternalServerError},
		{"BadRequest", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			}))
			defer server.Close()

			client := &clientImpl{
				oauthConfig: OauthConfig{
					ServerHost: server.URL,
				},
			}

			config := &request.Config{
				Host:    server.URL,
				Path:    "/validate",
				Headers: map[string]string{},
				Query:   map[string][]string{},
			}

			result, err := client.validate(config)

			assert.NoError(t, err)
			assert.False(t, result)
		})
	}
}

func TestClientImplCheckTokenHeadersAndQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Values("Authorization")[0]
		scope := r.URL.Query().Get("scope")
		hType := r.URL.Query().Get("type")
		assert.Equal(t, "Bearer token123", authorization)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "read", scope)
		assert.Equal(t, "access_token", hType)
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	client := &clientImpl{
		oauthConfig: OauthConfig{
			ServerHost: server.URL,
		},
	}

	input := Input{
		Headers: map[string]string{
			"Authorization": "Bearer token123",
			"Content-Type":  "application/json",
		},
		QueryParams: map[string][]string{
			"scope": {"read"},
			"type":  {"access_token"},
		},
	}

	result, err := client.CheckToken(input)

	assert.NoError(t, err)
	assert.True(t, result)
}
