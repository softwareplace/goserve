package server

import (
	"github.com/softwareplace/http-utils/security"
	"net/http"
)

func (a *apiRouterHandlerImpl) PublicRouter(handler ApiContextHandler, path string, method string) {
	a.Add(handler, path, method)
	security.AddOpenPath(method + "::" + ContextPath + path)
}

func (a *apiRouterHandlerImpl) Add(handler ApiContextHandler, path string, method string, requiredRoles ...string) {
	apiRoute.HandleFunc(ContextPath+path, func(writer http.ResponseWriter, req *http.Request) {
		ctx := Of(writer, req, "ROUTER/HANDLER")
		handler(&ctx)
	}).Methods(method)

	security.AddRoles(method+"::"+ContextPath+path, requiredRoles...)
}

func (a *apiRouterHandlerImpl) Get(handler ApiContextHandler, path string, requiredRoles ...string) {
	a.Add(handler, path, "GET", requiredRoles...)
}

func (a *apiRouterHandlerImpl) Post(handler ApiContextHandler, path string, requiredRoles ...string) {
	a.Add(handler, path, "POST", requiredRoles...)
}

func (a *apiRouterHandlerImpl) Put(handler ApiContextHandler, path string, requiredRoles ...string) {
	a.Add(handler, path, "PUT", requiredRoles...)
}

func (a *apiRouterHandlerImpl) Delete(handler ApiContextHandler, path string, requiredRoles ...string) {
	a.Add(handler, path, "DELETE", requiredRoles...)
}

func (a *apiRouterHandlerImpl) Patch(handler ApiContextHandler, path string, requiredRoles ...string) {
	a.Add(handler, path, "PATCH", requiredRoles...)
}

func (a *apiRouterHandlerImpl) Options(handler ApiContextHandler, path string, requiredRoles ...string) {
	a.Add(handler, path, "OPTIONS", requiredRoles...)
}

func (a *apiRouterHandlerImpl) Head(handler ApiContextHandler, path string, requiredRoles ...string) {
	a.Add(handler, path, "HEAD", requiredRoles...)
}
