package server

import (
	"github.com/gorilla/mux"
	"github.com/softwareplace/http-utils/api_context"
	"net/http"
)

type ApiContextHandler[T api_context.ApiContextData] func(ctx *api_context.ApiRequestContext[T])
type ApiMiddleware[T api_context.ApiContextData] func(*api_context.ApiRequestContext[T]) (doNext bool)

type ApiRouterHandler[T api_context.ApiContextData] interface {
	PublicRouter(handler ApiContextHandler[T], path string, method string)
	Add(handler ApiContextHandler[T], path string, method string, requiredRoles ...string)
	Get(handler ApiContextHandler[T], path string, requiredRoles ...string)
	Post(handler ApiContextHandler[T], path string, requiredRoles ...string)
	Put(handler ApiContextHandler[T], path string, requiredRoles ...string)
	Delete(handler ApiContextHandler[T], path string, requiredRoles ...string)
	Patch(handler ApiContextHandler[T], path string, requiredRoles ...string)
	Options(handler ApiContextHandler[T], path string, requiredRoles ...string)
	Head(handler ApiContextHandler[T], path string, requiredRoles ...string)
	StartServer()
	Use(middleware ApiMiddleware[T], name string)
}

type apiRouterHandlerImpl[T api_context.ApiContextData] struct {
	router *mux.Router
}

func (a *apiRouterHandlerImpl[T]) Use(middleware ApiMiddleware[T], name string) {
	a.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := api_context.Of[T](w, r, name)
			if middleware(ctx) {
				ctx.Next(next)
			}
		})
	})
}

func New[T api_context.ApiContextData](topMiddlewares ...ApiMiddleware[T]) ApiRouterHandler[T] {
	router := mux.NewRouter()
	router.Use(rootAppMiddleware[T])

	api := &apiRouterHandlerImpl[T]{
		router: router,
	}

	for _, middleware := range topMiddlewares {
		api.Use(middleware, "")
	}

	router.Use(HasResourceAccess[T])
	return api
}

func NewApiWith[T api_context.ApiContextData](router *mux.Router) ApiRouterHandler[T] {
	router.Use(rootAppMiddleware[T])
	router.Use(HasResourceAccess[T])

	api := &apiRouterHandlerImpl[T]{
		router: router,
	}

	return api
}
