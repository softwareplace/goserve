package impl

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/softwareplace/goserve/authorization/config"
	"github.com/softwareplace/goserve/authorization/input"
	authrequest "github.com/softwareplace/goserve/authorization/request"
	"github.com/softwareplace/goserve/authorization/response"
	authresponse "github.com/softwareplace/goserve/authorization/response"
	"github.com/softwareplace/goserve/request"
	"github.com/softwareplace/goserve/validator"
)

// ClientImpl struct
type ClientImpl struct {
	OauthConfig config.OauthConfig
}

// CheckToken checks if the token is valid
func (c *ClientImpl) CheckToken(input input.Input) (bool, error) {
	config := request.Build(c.OauthConfig.ServerHost).
		WithPath("authorization")

	for name, header := range input.Headers {
		if header == "" {
			continue
		}

		config.WithHeader(name, header)
	}

	for name, value := range input.QueryParams {
		if value == nil {
			continue
		}

		for _, v := range value {
			config.WithQuery(name, v)
		}
	}

	return c.validate(config)
}

// ChecktokenCustom checks if the token is valid
func (c *ClientImpl) CheckTokenCustom(config *request.Config) (bool, error) {
	return c.validate(config)
}

func (c *ClientImpl) validate(config *request.Config) (bool, error) {
	client := request.NewService()
	response, err := client.Get(config)

	defer client.Close()

	if response == nil {
		return false, fmt.Errorf("no response available")
	}

	if err != nil {
		return false, err
	}

	if response.StatusCode == http.StatusOK {
		return true, nil
	}

	log.Errorf("Failed to validate token with response status code: %v", response.StatusCode)

	return false, nil
}

func (c *ClientImpl) Login(
	authRequest authrequest.AuhtorizationRequest,
	applicationID string,
) (*response.AuthorizationResponse, error) {
	if err := validator.StructValidation(authRequest); err != nil {
		return nil, err
	}

	if len(applicationID) == 0 {
		return nil, fmt.Errorf("applicationID is required")
	}

	config := request.Build(c.OauthConfig.ServerHost).
		WithPath("login").
		WithHeader(request.RequestClientID, applicationID).
		WithBody(authRequest)

	client := request.NewService()

	defer client.Close()

	response, err := client.Post(config)

	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, fmt.Errorf("no response available")
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to login with response status code: %v", response.StatusCode)
	}

	responseData := authresponse.AuthorizationResponse{}

	err = client.BodyDecode(&responseData)

	if err != nil {
		return nil, err
	}

	if err := validator.StructValidation(responseData); err != nil {
		return nil, fmt.Errorf("failed to validate response data: %v", err)
	}

	return &responseData, nil
}
