package request

import "fmt"

type Config struct {
	Host               string
	Path               string
	Headers            map[string]string
	Query              map[string]string
	Body               any
	ExpectedStatusCode int
}

type Api[T any] interface {
	Get(config *Config) (*T, error)
	Post(config *Config) (*T, error)
	Put(config *Config) (*T, error)
	Delete(config *Config) (*T, error)
	Patch(config *Config) (*T, error)
	Head(config *Config) (*T, error)
}

func NewApi[T any](response T) Api[T] {
	i := new(_impl[T])
	i.response = response
	return i
}

type _impl[T any] struct {
	response T
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
