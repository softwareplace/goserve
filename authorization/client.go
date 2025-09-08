package authorization

import "github.com/softwareplace/goserve/request"

// Input struct
type Input struct {
	Headers     map[string]string // Headers
	QueryParams map[string][]string // Query parameters
}

type AuthorizationResponse struct {
	Jwt      string `json:"jwt" validate:"required,gt=20"`
	Expires  int64  `json:"expires" validate:"required,gt=1757353373"`
	IssuedAt int64  `json:"issuedAt" validate:"required,gt=1757353373"`
}

type AuhtorizationRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Client struct
type Client interface {
	CheckToken(input Input) (bool, error)
	ChecktokenCustom(config *request.Config) (bool, error)
	Login(request AuhtorizationRequest, applicationID string) (*AuthorizationResponse, error)
}

// NewClient creates a new client
func NewClient(oauthConfig OauthConfig) Client {
	return &clientImpl{
		oauthConfig: oauthConfig,
	}
}
