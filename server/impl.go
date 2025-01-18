package server

import (
	"github.com/gorilla/mux"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/error_handler"
	"github.com/softwareplace/http-utils/security/principal"
	"net/http"
)

type apiRouterHandlerImpl[T api_context.ApiPrincipalContext] struct {
	router           *mux.Router
	principalService *principal.PService[T]
	errorHandler     *error_handler.ApiErrorHandler[T]
	loginService     *LoginService[T]
}

func (a *apiRouterHandlerImpl[T]) RegisterMiddleware(middleware ApiMiddleware[T], name string) ApiRouterHandler[T] {
	a.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := api_context.Of[T](w, r, name)
			if middleware(ctx) {
				ctx.Next(next)
			}
		})
	})
	return a
}

func (a *apiRouterHandlerImpl[T]) WithErrorHandler(handler *error_handler.ApiErrorHandler[T]) ApiRouterHandler[T] {
	a.errorHandler = handler
	return a
}

func (a *apiRouterHandlerImpl[T]) WithLoginResource(loginService *LoginService[T]) ApiRouterHandler[T] {
	a.loginService = loginService
	a.PublicRouter(a.Login, "login", "POST")
	return a
}

func CreateApiRouter[T api_context.ApiPrincipalContext](topMiddlewares ...ApiMiddleware[T]) ApiRouterHandler[T] {
	router := mux.NewRouter()
	router.Use(rootAppMiddleware[T])

	api := &apiRouterHandlerImpl[T]{
		router: router,
	}

	router.Use(api.errorHandlerWrapper)

	for _, middleware := range topMiddlewares {
		api.RegisterMiddleware(middleware, "")
	}
	return api
}

func (a *apiRouterHandlerImpl[T]) WithPrincipalService(service *principal.PService[T]) ApiRouterHandler[T] {
	a.principalService = service
	if a.principalService != nil {
		a.router.Use(a.hasResourceAccess)
	}
	return a
}

func CreateApiRouterWith[T api_context.ApiPrincipalContext](router *mux.Router) ApiRouterHandler[T] {
	router.Use(rootAppMiddleware[T])
	api := &apiRouterHandlerImpl[T]{
		router: router,
	}

	return api
}
