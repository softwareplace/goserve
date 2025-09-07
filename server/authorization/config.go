package authorization

import "github.com/softwareplace/goserve/utils"

// OauthConfig struct
type OauthConfig struct {
	// ServerHost is the server host
	ServerHost string
}

// NewOauthConfig creates a new oauth config
func NewOauthConfig() OauthConfig {
	return OauthConfig{
		ServerHost: utils.GetRequiredEnv("OAUTH_SERVER_HOST"),
	}
}
