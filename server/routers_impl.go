package server

import (
	"net/http"
	"strings"

	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security/router"
)

func (a *baseServer[T]) PublicRouter(handler ApiContextHandler[T], path string, method string) Api[T] {
	a.Add(handler, path, method)
	combinedKey := method + "::" + a.contextPath + path
	router.AddOpenPath(combinedKey)
	return a
}

func (a *baseServer[T]) Add(handler ApiContextHandler[T], path string, method string, requiredRoles ...string) Api[T] {

	handlerPath := strings.TrimSuffix(a.contextPath, "/") + "/" + strings.TrimPrefix(path, "/")

	a.router.HandleFunc(handlerPath, func(writer http.ResponseWriter, req *http.Request) {
		ctx := goservectx.Of[T](writer, req, "ROUTER/HANDLER")
		handler(ctx)
	}).Methods(method)

	router.AddRoles(method+"::"+handlerPath, requiredRoles...)
	return a
}

func (a *baseServer[T]) Get(handler ApiContextHandler[T], path string, requiredRoles ...string) Api[T] {
	a.Add(handler, path, "GET", requiredRoles...)
	return a
}

func (a *baseServer[T]) Post(handler ApiContextHandler[T], path string, requiredRoles ...string) Api[T] {
	a.Add(handler, path, "POST", requiredRoles...)
	return a
}

func (a *baseServer[T]) Put(handler ApiContextHandler[T], path string, requiredRoles ...string) Api[T] {
	a.Add(handler, path, "PUT", requiredRoles...)
	return a
}

func (a *baseServer[T]) Delete(handler ApiContextHandler[T], path string, requiredRoles ...string) Api[T] {
	a.Add(handler, path, "DELETE", requiredRoles...)
	return a
}

func (a *baseServer[T]) Patch(handler ApiContextHandler[T], path string, requiredRoles ...string) Api[T] {
	a.Add(handler, path, "PATCH", requiredRoles...)
	return a
}

func (a *baseServer[T]) Options(handler ApiContextHandler[T], path string, requiredRoles ...string) Api[T] {
	a.Add(handler, path, "OPTIONS", requiredRoles...)
	return a
}

func (a *baseServer[T]) Head(handler ApiContextHandler[T], path string, requiredRoles ...string) Api[T] {
	a.Add(handler, path, "HEAD", requiredRoles...)
	return a
}
