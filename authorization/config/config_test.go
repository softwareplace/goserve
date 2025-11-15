package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOauthConfigWithEnvVar(t *testing.T) {
	expectedHost := "https://oauth.test.com"

	t.Setenv("OAUTH_SERVER_HOST", expectedHost)

	config := NewOauthConfig()

	assert.Equal(t, expectedHost, config.ServerHost)
}

func TestNewOauthConfigWithoutEnvVarShouldPanic(t *testing.T) {
	assert.Panics(t, func() {
		NewOauthConfig()
	}, "NewOauthConfig should panic when OAUTH_SERVER_HOST is not set")
}

func TestNewOauthConfigWithEmptyEnvVarShouldPanic(t *testing.T) {
	t.Setenv("OAUTH_SERVER_HOST", "")

	assert.Panics(t, func() {
		NewOauthConfig()
	}, "NewOauthConfig should panic when OAUTH_SERVER_HOST is empty")
}

func TestOauthConfigStructure(t *testing.T) {
	expectedHost := "https://oauth.example.com"

	config := OauthConfig{
		ServerHost: expectedHost,
	}

	assert.Equal(t, expectedHost, config.ServerHost)
}

func TestNewOauthConfigWithDifferentHostFormats(t *testing.T) {
	testCases := []struct {
		name         string
		hostValue    string
		expectedHost string
	}{
		{
			name:         "HTTPS URL",
			hostValue:    "https://oauth.example.com",
			expectedHost: "https://oauth.example.com",
		},
		{
			name:         "HTTP URL",
			hostValue:    "http://oauth.example.com",
			expectedHost: "http://oauth.example.com",
		},
		{
			name:         "URL with port",
			hostValue:    "https://oauth.example.com:8080",
			expectedHost: "https://oauth.example.com:8080",
		},
		{
			name:         "URL with path",
			hostValue:    "https://oauth.example.com/auth",
			expectedHost: "https://oauth.example.com/auth",
		},
		{
			name:         "Localhost",
			hostValue:    "http://localhost:3000",
			expectedHost: "http://localhost:3000",
		},
		{
			name:         "IP address",
			hostValue:    "http://192.168.1.100:8080",
			expectedHost: "http://192.168.1.100:8080",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("OAUTH_SERVER_HOST", tc.hostValue)

			config := NewOauthConfig()

			assert.Equal(t, tc.expectedHost, config.ServerHost)
		})
	}
}
