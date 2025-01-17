package server

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/security/validator"
	"net/http"
)

func (a *apiRouterHandlerImpl[T]) PublicRouter(handler ApiContextHandler[T], path string, method string) {
	a.Add(handler, path, method)
	validator.AddOpenPath(method + "::" + ContextPath + path)
}

func (a *apiRouterHandlerImpl[T]) Add(handler ApiContextHandler[T], path string, method string, requiredRoles ...string) {
	a.router.HandleFunc(ContextPath+path, func(writer http.ResponseWriter, req *http.Request) {
		ctx := api_context.Of[T](writer, req, "ROUTER/HANDLER")
		handler(ctx)
	}).Methods(method)

	validator.AddRoles(method+"::"+ContextPath+path, requiredRoles...)
}

func (a *apiRouterHandlerImpl[T]) Get(handler ApiContextHandler[T], path string, requiredRoles ...string) {
	a.Add(handler, path, "GET", requiredRoles...)
}

func (a *apiRouterHandlerImpl[T]) Post(handler ApiContextHandler[T], path string, requiredRoles ...string) {
	a.Add(handler, path, "POST", requiredRoles...)
}

func (a *apiRouterHandlerImpl[T]) Put(handler ApiContextHandler[T], path string, requiredRoles ...string) {
	a.Add(handler, path, "PUT", requiredRoles...)
}

func (a *apiRouterHandlerImpl[T]) Delete(handler ApiContextHandler[T], path string, requiredRoles ...string) {
	a.Add(handler, path, "DELETE", requiredRoles...)
}

func (a *apiRouterHandlerImpl[T]) Patch(handler ApiContextHandler[T], path string, requiredRoles ...string) {
	a.Add(handler, path, "PATCH", requiredRoles...)
}

func (a *apiRouterHandlerImpl[T]) Options(handler ApiContextHandler[T], path string, requiredRoles ...string) {
	a.Add(handler, path, "OPTIONS", requiredRoles...)
}

func (a *apiRouterHandlerImpl[T]) Head(handler ApiContextHandler[T], path string, requiredRoles ...string) {
	a.Add(handler, path, "HEAD", requiredRoles...)
}
