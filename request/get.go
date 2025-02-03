package request

import "net/http"

func (i *_impl) Get(config *Config) (*http.Response, error) {
	return i.exec("GET", config)
}
