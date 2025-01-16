package server

import (
	"github.com/gorilla/mux"
	"github.com/softwareplace/http-utils/api_context"
)

type DefaultContext struct{}

func Default() ApiRouterHandler[DefaultContext] {
	api := &apiRouterHandlerImpl[DefaultContext]{
		router: mux.NewRouter(),
	}
	api.router.Use(rootAppMiddleware)
	return api
}

func (d DefaultContext) SetAuthorizationClaims(map[string]interface{}) {

}

func (d DefaultContext) SetApiKeyId(string) {

}

func (d DefaultContext) SetAccessId(string) {

}

func (d DefaultContext) Data(api_context.ApiContextData) {

}

func (d DefaultContext) Salt() string {
	return ""
}

func (d DefaultContext) Roles() []string {
	return []string{}
}
