package authorization

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/softwareplace/goserve/request"
)

// clientImpl struct
type clientImpl struct {
	oauthConfig OauthConfig
}

// CheckToken checks if the token is valid
func (c *clientImpl) CheckToken(input Input) (bool, error) {
	config := request.Build(c.oauthConfig.ServerHost)

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
func (c *clientImpl) ChecktokenCustom(config *request.Config) (bool, error) {
	return c.validate(config)
}

func (c *clientImpl) validate(config *request.Config) (bool, error) {
	client := request.NewService()
	response, err := client.Get(config)

	go func() {
		client.Close()
	}()

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
