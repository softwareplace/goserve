package request

import "net/http"

func (i *_impl) Patch(config *Config) (*http.Response, error) {
	return i.exec("PATCH", config)
}
