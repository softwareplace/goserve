package authorization

import "github.com/softwareplace/goserve/env"

// OauthConfig struct
type OauthConfig struct {
	// ServerHost is the server host
	ServerHost string
}

// NewOauthConfig creates a new oauth config
func NewOauthConfig() OauthConfig {
	return OauthConfig{
		ServerHost: env.GetRequiredEnv("OAUTH_SERVER_HOST"),
	}
}
