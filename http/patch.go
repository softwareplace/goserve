package http

import "net/http"

func (i *_impl) Patch(config *Config) (*http.Response, error) {
	return i.Exec("PATCH", config)
}
