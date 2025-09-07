package authorization

import "github.com/softwareplace/goserve/request"

// Input struct
type Input struct {
	Headers     map[string]string // Headers
	QueryParams map[string][]string // Query parameters
}

// Client struct
type Client interface {
	CheckToken(input Input) (bool, error)
	ChecktokenCustom(config *request.Config) (bool, error)
}

// NewClient creates a new client
func NewClient(oauthConfig OauthConfig) Client {
	return &clientImpl{
		oauthConfig: oauthConfig,
	}
}
