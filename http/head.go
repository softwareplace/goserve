package http

import "net/http"

func (i *_impl) Head(config *Config) (*http.Response, error) {
	return i.Exec("HEAD", config)
}
