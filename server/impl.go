package server

import (
	"github.com/gorilla/mux"
	apicontext "github.com/softwareplace/http-utils/context"
	errorhandler "github.com/softwareplace/http-utils/error"
	"github.com/softwareplace/http-utils/security/principal"
	"net/http"
	"strings"
)

func (a *apiRouterHandlerImpl[T]) RegisterMiddleware(middleware ApiMiddleware[T], name string) Api[T] {
	a.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := apicontext.Of[T](w, r, name)
			if middleware(ctx) {
				ctx.Next(next)
			}
		})
	})
	return a
}

func (a *apiRouterHandlerImpl[T]) RegisterCustomMiddleware(middleware func(next http.Handler) http.Handler) Api[T] {
	a.router.Use(middleware)
	return a
}

func (a *apiRouterHandlerImpl[T]) ErrorHandler(handler errorhandler.ApiHandler[T]) Api[T] {
	a.errorHandler = handler
	return a
}

func (a *apiRouterHandlerImpl[T]) LoginResource(loginService LoginService[T]) Api[T] {
	a.loginService = loginService
	if a.loginResourceEnable {
		a.PublicRouter(a.Login, "login", "POST")
	}
	return a
}

func (a *apiRouterHandlerImpl[T]) ApiKeyGeneratorResource(apiKeyGeneratorService ApiKeyGeneratorService[T]) Api[T] {
	a.apiKeyGeneratorService = apiKeyGeneratorService
	if a.apiSecretKeyGeneratorResourceEnable {
		a.Post(a.ApiKeyGenerator, "api-key/generate", "POST", strings.Join(apiKeyGeneratorService.RequiredScopes(), " "))
	}
	return a
}

func (a *apiRouterHandlerImpl[T]) ApiSecretKeyGeneratorResourceEnabled(enable bool) Api[T] {
	a.apiSecretKeyGeneratorResourceEnable = enable
	return a
}

func (a *apiRouterHandlerImpl[T]) LoginResourceEnabled(enable bool) Api[T] {
	a.loginResourceEnable = enable
	return a
}

func (a *apiRouterHandlerImpl[T]) PrincipalService(service principal.Service[T]) Api[T] {
	a.principalService = service
	if a.principalService != nil {
		a.router.Use(a.hasResourceAccess)
	}
	return a
}

func (a *apiRouterHandlerImpl[T]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.router.ServeHTTP(w, req)
}

func (a *apiRouterHandlerImpl[T]) Router() *mux.Router {
	return a.router
}

func (a *apiRouterHandlerImpl[T]) RouterHandler(handler RouterHandler) Api[T] {
	handler(a.router)
	return a
}

func (a *apiRouterHandlerImpl[T]) EmbeddedServer(handler func(Api[T])) Api[T] {
	handler(a)
	return a
}
