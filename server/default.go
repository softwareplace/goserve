package server

import (
	"github.com/gorilla/mux"
	"github.com/softwareplace/http-utils/api_context"
)

func Default() ApiRouterHandler[*api_context.DefaultContext] {
	api := &apiRouterHandlerImpl[*api_context.DefaultContext]{
		router: mux.NewRouter(),
	}
	api.router.Use(rootAppMiddleware)
	return api
}
