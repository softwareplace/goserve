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
	Router() *mux.Router
	StartServer()
}

type apiRouterHandlerImpl struct {
	router *mux.Router
}

func (a *apiRouterHandlerImpl) Router() *mux.Router {
	return a.router
}

func New() ApiRouterHandler {
	api := &apiRouterHandlerImpl{
		router: mux.NewRouter(),
	}
	api.router.Use(rootAppMiddleware)
	return api
}

func NewApiWith(router *mux.Router) ApiRouterHandler {
	api := &apiRouterHandlerImpl{
		router: router,
	}
	api.router.Use(rootAppMiddleware)
	return api
}
