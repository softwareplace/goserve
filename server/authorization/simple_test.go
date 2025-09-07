package authorization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleOauthConfig(t *testing.T) {
	t.Setenv("OAUTH_SERVER_HOST", "https://test.example.com")

	config := NewOauthConfig()
	assert.Equal(t, "https://test.example.com", config.ServerHost)
}

func TestSimpleNewClient(t *testing.T) {
	oauthConfig := OauthConfig{
		ServerHost: "https://oauth.example.com",
	}

	client := NewClient(oauthConfig)
	assert.NotNil(t, client)
}
