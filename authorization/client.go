package authorization

import (
	"github.com/softwareplace/goserve/authorization/config"
	"github.com/softwareplace/goserve/authorization/input"
	"github.com/softwareplace/goserve/authorization/internal/impl"
	authrequest "github.com/softwareplace/goserve/authorization/request"
	"github.com/softwareplace/goserve/authorization/response"
	"github.com/softwareplace/goserve/request"
)

// Client struct
type Client interface {
	CheckToken(input input.Input) (bool, error)
	CheckTokenCustom(config *request.Config) (bool, error)
	Login(request authrequest.AuhtorizationRequest, applicationID string) (*response.AuthorizationResponse, error)
}

// NewClient creates a new client
func NewClient(oauthConfig config.OauthConfig) Client {
	return &impl.ClientImpl{
		OauthConfig: oauthConfig,
	}
}
