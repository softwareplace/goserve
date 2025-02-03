package request

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Config struct {
	Host               string
	Path               string
	Headers            map[string]string
	Query              map[string]string
	Body               any
	ExpectedStatusCode int
}

type Service interface {
	// Get sends an HTTP GET request with the provided configuration and returns the HTTP response or an error.
	Get(config *Config) (*http.Response, error)

	// Post sends an HTTP POST request with the provided configuration and returns the HTTP response or an error.
	Post(config *Config) (*http.Response, error)

	// Put sends an HTTP PUT request with the provided configuration and returns the HTTP response or an error.
	Put(config *Config) (*http.Response, error)

	// Delete sends an HTTP DELETE request with the provided configuration and returns the HTTP response or an error.
	Delete(config *Config) (*http.Response, error)

	// Patch sends an HTTP PATCH request with the provided configuration and returns the HTTP response or an error.
	Patch(config *Config) (*http.Response, error)

	// Head sends an HTTP HEAD request with the provided configuration and returns the HTTP response or an error.
	Head(config *Config) (*http.Response, error)

	// ToString converts the body of the last HTTP response into a string and returns it along with an error if any.
	ToString() (string, error)

	// BodyDecode decodes the body of the last HTTP response into the given target interface and returns an error if any.
	// Example:
	//   ...
	//	 data := MyData{}
	//   err := client.BodyDecode(&data)
	BodyDecode(target any) error

	// Close closes the response body of the last HTTP response to release resources.
	Close()
}

func NewService() Service {
	i := new(_impl)
	return i
}

type _impl struct {
	response *http.Response
}

func (i *_impl) Close() {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(i.response.Body)
}

func (i *_impl) ToString() (string, error) {

	// Ensure the response is not nil
	if i.response == nil {
		return "", fmt.Errorf("no response available")
	}

	// Read the response body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(i.response.Body)

	bodyBytes, err := io.ReadAll(i.response.Body)

	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(bodyBytes), nil
}

func (i *_impl) BodyDecode(target any) error {

	if i.response == nil {
		return fmt.Errorf("no response available")
	}

	decoder := json.NewDecoder(i.response.Body)
	err := decoder.Decode(target)
	if err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}

	return nil
}

func Build(host string) *Config {
	config := &Config{}
	config.Host = host
	config.Path = ""
	config.Headers = map[string]string{}
	config.Query = map[string]string{}
	config.Body = nil
	config.ExpectedStatusCode = 200
	config.WithHeader("Content-Type", "application/json")
	return config
}

func (config *Config) WithPath(path string) *Config {
	config.Path = path
	return config
}

func (config *Config) WithQuery(name string, value any) *Config {
	config.Query[name] = fmt.Sprintf("%v", value)
	return config
}

func (config *Config) WithHeader(name string, value any) *Config {
	config.Headers[name] = fmt.Sprintf("%v", value)
	return config
}

func (config *Config) WithBody(body any) *Config {
	config.Body = body
	return config
}

func (config *Config) WithExpectedStatusCode(expectedStatusCode int) *Config {
	config.ExpectedStatusCode = expectedStatusCode
	return config
}
