package request

import "net/http"

func (i *_impl) Delete(config *Config) (*http.Response, error) {
	return i.exec("DELETE", config)
}
