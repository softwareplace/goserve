package request

import "net/http"

func (i *_impl) Post(config *Config) (*http.Response, error) {
	return i.exec("POST", config)
}
