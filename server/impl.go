package server

import (
	"github.com/gorilla/mux"
	"github.com/softwareplace/http-utils/apicontext"
	"github.com/softwareplace/http-utils/error_handler"
	"github.com/softwareplace/http-utils/security/principal"
	"net/http"
	"strings"
)

func (a *apiRouterHandlerImpl[T]) RegisterMiddleware(middleware ApiMiddleware[T], name string) ApiRouterHandler[T] {
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

func (a *apiRouterHandlerImpl[T]) RegisterCustomMiddleware(middleware func(next http.Handler) http.Handler) ApiRouterHandler[T] {
	a.router.Use(middleware)
	return a
}

func (a *apiRouterHandlerImpl[T]) WithErrorHandler(handler error_handler.ApiErrorHandler[T]) ApiRouterHandler[T] {
	a.errorHandler = handler
	return a
}

func (a *apiRouterHandlerImpl[T]) WithLoginResource(loginService LoginService[T]) ApiRouterHandler[T] {
	a.loginService = loginService
	if a.loginResourceEnable {
		a.PublicRouter(a.Login, "login", "POST")
	}
	return a
}

func (a *apiRouterHandlerImpl[T]) WithApiKeyGeneratorResource(apiKeyGeneratorService ApiKeyGeneratorService[T]) ApiRouterHandler[T] {
	a.apiKeyGeneratorService = apiKeyGeneratorService
	if a.apiSecretKeyGeneratorResourceEnable {
		a.Post(a.ApiKeyGenerator, "api-key/generate", "POST", strings.Join(apiKeyGeneratorService.RequiredScopes(), " "))
	}
	return a
}

func (a *apiRouterHandlerImpl[T]) ApiSecretKeyGeneratorResourceEnabled(enable bool) ApiRouterHandler[T] {
	a.apiSecretKeyGeneratorResourceEnable = enable
	return a
}

func (a *apiRouterHandlerImpl[T]) LoginResourceEnabled(enable bool) ApiRouterHandler[T] {
	a.loginResourceEnable = enable
	return a
}

func (a *apiRouterHandlerImpl[T]) WithPrincipalService(service principal.PService[T]) ApiRouterHandler[T] {
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

func (a *apiRouterHandlerImpl[T]) RouterHandler(handler RouterHandler) ApiRouterHandler[T] {
	handler(a.router)
	return a
}

func (a *apiRouterHandlerImpl[T]) EmbeddedServer(handler func(ApiRouterHandler[T])) ApiRouterHandler[T] {
	handler(a)
	return a
}
