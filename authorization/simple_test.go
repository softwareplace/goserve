package authorization

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/softwareplace/goserve/authorization/config"
)

func TestSimpleOauthConfig(t *testing.T) {
	t.Setenv("OAUTH_SERVER_HOST", "https://test.example.com")

	config := config.NewOauthConfig()
	assert.Equal(t, "https://test.example.com", config.ServerHost)
}

func TestSimpleNewClient(t *testing.T) {
	oauthConfig := config.OauthConfig{
		ServerHost: "https://oauth.example.com",
	}

	client := NewClient(oauthConfig)
	assert.NotNil(t, client)
}
