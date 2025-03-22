package server

import (
	"github.com/gorilla/mux"
	"github.com/softwareplace/http-utils/apicontext"
)

func CreateApiRouter[T apicontext.ApiPrincipalContext](topMiddlewares ...ApiMiddleware[T]) ApiRouterHandler[T] {
	router := mux.NewRouter()
	router.Use(rootAppMiddleware[T])

	api := &apiRouterHandlerImpl[T]{
		router:                              router,
		apiSecretKeyGeneratorResourceEnable: true,
		loginResourceEnable:                 true,
		contextPath:                         apiContextPath(),
		port:                                apiPort(),
	}

	router.Use(api.errorHandlerWrapper)

	for _, middleware := range topMiddlewares {
		api.RegisterMiddleware(middleware, "")
	}
	return api.NotFoundHandler()
}

func CreateApiRouterWith[T apicontext.ApiPrincipalContext](router mux.Router) ApiRouterHandler[T] {
	router.Use(rootAppMiddleware[T])
	api := &apiRouterHandlerImpl[T]{
		router:                              &router,
		apiSecretKeyGeneratorResourceEnable: true,
		loginResourceEnable:                 true,
		contextPath:                         apiContextPath(),
		port:                                apiPort(),
	}

	return api.NotFoundHandler()
}
