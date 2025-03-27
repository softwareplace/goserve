package server

import (
	"github.com/gorilla/mux"
	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/security"
	"github.com/softwareplace/goserve/security/login"
	"github.com/softwareplace/goserve/security/secret"
	"net/http"
	"sync"
)

type baseServer[T goservectx.Principal] struct {
	router                              *mux.Router
	errorHandler                        goserveerror.ApiHandler[T]
	loginService                        login.Service[T]
	securityService                     security.Service[T]
	secretService                       secret.Service[T]
	server                              *http.Server // Add a server instance
	mu                                  sync.Mutex   // Add a mutex for thread safety
	swaggerIsEnabled                    bool
	loginResourceEnable                 bool
	apiSecretKeyGeneratorResourceEnable bool
	contextPath                         string
	port                                string
}

// Default initializes and returns a new API instance configured to work with the DefaultContext type.
// It sets up the router, applies any provided top-level middlewares, and assigns default options
// such as the context path and port. This is intended for testing and development environments only,
// as DefaultContext is not secure for production use.
//
// Parameters:
//   - topMiddlewares: Optional list of middlewares that will be applied globally to the API.
//
// Returns:
//   - Api[*goservectx.DefaultContext]: An API instance configured with DefaultContext.
func Default(
	topMiddlewares ...ApiMiddleware[*goservectx.DefaultContext],
) Api[*goservectx.DefaultContext] {
	return New[*goservectx.DefaultContext](topMiddlewares...)
}

// New initializes and returns a new instance of the Api[T] interface.
// It sets up the router, adds the root application middleware and any provided
// top-level middlewares, and configures default options such as the context path and port.
//
// Parameters:
//   - T: A type that implements the goservectx.Principal interface.
//   - topMiddlewares: Optional list of middlewares to apply at the API level.
//
// Returns:
//   - Api[T]: An instance of the Api[T] interface with the configured router and default behaviors.
func New[T goservectx.Principal](topMiddlewares ...ApiMiddleware[T]) Api[T] {
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

// NewWith initializes and returns a new instance of the Api[T] interface using a provided Gorilla mux router.
// It wraps the provided router with the baseServer configuration, adds the root application middleware,
// and sets default options such as the context path and port.
//
// Parameters:
//   - T: A type that implements the goservectx.Principal interface.
//   - router: An instance of mux.Router to be configured and used by the API.
//
// Returns:
//   - Api[T]: An instance of the Api[T] interface configured with the provided router.
func NewWith[T goservectx.Principal](router mux.Router) Api[T] {
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
