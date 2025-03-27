package server

import (
	"github.com/gorilla/mux"
	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security/login"
	"net/http"
)

func (a *baseServer[T]) RegisterCustomMiddleware(middleware func(next http.Handler) http.Handler) Api[T] {
	a.router.Use(middleware)
	return a
}

func (a *baseServer[T]) ErrorHandler(handler goservectx.ApiHandler[T]) Api[T] {
	a.errorHandler = handler
	return a
}

func (a *baseServer[T]) LoginService(service login.Service[T]) Api[T] {
	a.loginService = service
	if a.loginResourceEnable {
		a.PublicRouter(a.Login, "login", "POST")
	}
	return a
}

func (a *baseServer[T]) SecretKeyGeneratorResourceEnabled(enable bool) Api[T] {
	a.apiSecretKeyGeneratorResourceEnable = enable
	return a
}

func (a *baseServer[T]) LoginResourceEnabled(enable bool) Api[T] {
	a.loginResourceEnable = enable
	return a
}

func (a *baseServer[T]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.router.ServeHTTP(w, req)
}

func (a *baseServer[T]) Router() *mux.Router {
	return a.router
}

func (a *baseServer[T]) RouterHandler(handler RouterHandler) Api[T] {
	handler(a.router)
	return a
}

func (a *baseServer[T]) EmbeddedServer(handler func(Api[T])) Api[T] {
	handler(a)
	return a
}
