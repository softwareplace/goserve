package request

import "net/http"

func (i *_impl) Put(config *Config) (*http.Response, error) {
	return i.exec("PUT", config)
}
