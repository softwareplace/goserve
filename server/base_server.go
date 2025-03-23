package server

import (
	"github.com/gorilla/mux"
	apicontext "github.com/softwareplace/http-utils/context"
	errorhandler "github.com/softwareplace/http-utils/error"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/security/principal"
	"github.com/softwareplace/http-utils/security/secret"
	"net/http"
	"sync"
)

type baseServer[T apicontext.Principal] struct {
	router                              *mux.Router
	principalService                    principal.Service[T]
	errorHandler                        errorhandler.ApiHandler[T]
	loginService                        LoginService[T]
	securityService                     security.Service[T]
	secretService                       secret.Service[T]
	apiKeyGeneratorService              ApiKeyGeneratorService[T]
	server                              *http.Server // Add a server instance
	mu                                  sync.Mutex   // Add a mutex for thread safety
	swaggerIsEnabled                    bool
	loginResourceEnable                 bool
	apiSecretKeyGeneratorResourceEnable bool
	contextPath                         string
	port                                string
}

func Default(
	topMiddlewares ...ApiMiddleware[*apicontext.DefaultContext],
) Api[*apicontext.DefaultContext] {
	return CreateApiRouter[*apicontext.DefaultContext](topMiddlewares...)
}

func CreateApiRouter[T apicontext.Principal](topMiddlewares ...ApiMiddleware[T]) Api[T] {
	router := mux.NewRouter()
	router.Use(rootAppMiddleware[T])

	api := &baseServer[T]{
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

func CreateApiRouterWith[T apicontext.Principal](router mux.Router) Api[T] {
	router.Use(rootAppMiddleware[T])
	api := &baseServer[T]{
		router:                              &router,
		apiSecretKeyGeneratorResourceEnable: true,
		loginResourceEnable:                 true,
		contextPath:                         apiContextPath(),
		port:                                apiPort(),
	}

	return api.NotFoundHandler()
}
