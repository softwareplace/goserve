package server

import (
	"github.com/gorilla/mux"
)

type ApiContextHandler func(ctx *ApiRequestContext)

type ApiRouterHandler interface {
	PublicRouter(handler ApiContextHandler, path string, method string)
	Add(handler ApiContextHandler, path string, method string, requiredRoles ...string)
	Get(handler ApiContextHandler, path string, requiredRoles ...string)
	Post(handler ApiContextHandler, path string, requiredRoles ...string)
	Put(handler ApiContextHandler, path string, requiredRoles ...string)
	Delete(handler ApiContextHandler, path string, requiredRoles ...string)
	Patch(handler ApiContextHandler, path string, requiredRoles ...string)
	Options(handler ApiContextHandler, path string, requiredRoles ...string)
	Head(handler ApiContextHandler, path string, requiredRoles ...string)
	StartServer()
}

type apiRouterHandlerImpl struct {
	Router *mux.Router
}

func New() ApiRouterHandler {
	api := &apiRouterHandlerImpl{
		Router: mux.NewRouter(),
	}
	api.Router.Use(rootAppMiddleware)
	return api
}

func NewApiWith(router *mux.Router) ApiRouterHandler {
	api := &apiRouterHandlerImpl{
		Router: router,
	}
	api.Router.Use(rootAppMiddleware)
	return api
}
