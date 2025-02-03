package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func (i *_impl) build(method string, config *Config) (*http.Request, error) {
	var body io.Reader

	if config.Body != nil {
		jsonBody, err := json.Marshal(config.Body)
		if err != nil {
			log.Fatalf("Failed to marshal request body: %v", err)
		}
		body = bytes.NewBuffer(jsonBody)
	}

	requestHost := strings.Trim(config.Host, "/")
	requestPath := strings.TrimPrefix(config.Path, "/")
	req, err := http.NewRequest(method, requestHost+"/"+requestPath, body)

	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	// Add headers to the request
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// Add query parameters to the URL
	query := req.URL.Query()
	for key, value := range config.Query {
		query.Set(key, value)
	}

	return req, nil
}

func (i *_impl) exec(method string, config *Config) (*http.Response, error) {
	request, err := i.build(method, config)

	if err != nil {
		return nil, fmt.Errorf("failed to build request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(request)

	i.response = resp
	return resp, nil
}
