package server

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/security/principal"
	"net/http"
	"strings"
)

func (a *apiRouterHandlerImpl[T]) PublicRouter(handler ApiContextHandler[T], path string, method string) ApiRouterHandler[T] {
	a.Add(handler, path, method)
	combinedKey := method + "::" + ContextPath + path
	principal.AddOpenPath(combinedKey)
	return a
}

func (a *apiRouterHandlerImpl[T]) Add(handler ApiContextHandler[T], path string, method string, requiredRoles ...string) ApiRouterHandler[T] {
	handlerPath := strings.TrimSuffix(ContextPath, "/") + "/" + strings.TrimPrefix(path, "/")

	a.router.HandleFunc(handlerPath, func(writer http.ResponseWriter, req *http.Request) {
		ctx := api_context.Of[T](writer, req, "ROUTER/HANDLER")
		handler(ctx)
	}).Methods(method)

	principal.AddRoles(method+"::"+handlerPath, requiredRoles...)
	return a
}

func (a *apiRouterHandlerImpl[T]) Get(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T] {
	a.Add(handler, path, "GET", requiredRoles...)
	return a
}

func (a *apiRouterHandlerImpl[T]) Post(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T] {
	a.Add(handler, path, "POST", requiredRoles...)
	return a
}

func (a *apiRouterHandlerImpl[T]) Put(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T] {
	a.Add(handler, path, "PUT", requiredRoles...)
	return a
}

func (a *apiRouterHandlerImpl[T]) Delete(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T] {
	a.Add(handler, path, "DELETE", requiredRoles...)
	return a
}

func (a *apiRouterHandlerImpl[T]) Patch(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T] {
	a.Add(handler, path, "PATCH", requiredRoles...)
	return a
}

func (a *apiRouterHandlerImpl[T]) Options(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T] {
	a.Add(handler, path, "OPTIONS", requiredRoles...)
	return a
}

func (a *apiRouterHandlerImpl[T]) Head(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T] {
	a.Add(handler, path, "HEAD", requiredRoles...)
	return a
}
