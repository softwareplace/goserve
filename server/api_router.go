package server

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/security/principal"
)

type ApiContextHandler[T api_context.ApiContextData] func(ctx *api_context.ApiRequestContext[T])

type ApiMiddleware[T api_context.ApiContextData] func(*api_context.ApiRequestContext[T]) (doNext bool)

type ApiRouterHandler[T api_context.ApiContextData] interface {
	// PublicRouter registers a public route handler that does not require authentication or authorization.
	// It allows unrestricted access from any client.
	//
	// Parameters:
	// - handler: The handler function to process requests.
	// - path: The URL route path.
	// - method: The HTTP method for the route.
	//
	// Returns:
	// - ApiRouterHandler: The router handler for chaining further route configurations.
	PublicRouter(handler ApiContextHandler[T], path string, method string) ApiRouterHandler[T]

	// Add registers a route handler with optional role-based access control.
	// This method is used to define routes and assign roles required to access them.
	//
	// Parameters:
	// - handler: The handler function to process requests.
	// - path: The URL route path.
	// - method: The HTTP method for the route.
	// - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	// - ApiRouterHandler: The router handler for chaining further route configurations.
	Add(handler ApiContextHandler[T], path string, method string, requiredRoles ...string) ApiRouterHandler[T]

	// Get registers a route handler specifically for HTTP GET requests with optional role-based access control.
	//
	// Parameters:
	// - handler: The handler function to process requests.
	// - path: The URL route path.
	// - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	// - ApiRouterHandler: The router handler for chaining further route configurations.
	Get(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T]

	// Post registers a route handler specifically for HTTP POST requests with optional role-based access control.
	//
	// Parameters:
	// - handler: The handler function to process requests.
	// - path: The URL route path.
	// - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	// - ApiRouterHandler: The router handler for chaining further route configurations.
	Post(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T]

	// Put registers a route handler specifically for HTTP PUT requests with optional role-based access control.
	//
	// Parameters:
	// - handler: The handler function to process requests.
	// - path: The URL route path.
	// - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	// - ApiRouterHandler: The router handler for chaining further route configurations.
	Put(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T]

	// Delete registers a route handler specifically for HTTP DELETE requests with optional role-based access control.
	//
	// Parameters:
	// - handler: The handler function to process requests.
	// - path: The URL route path.
	// - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	// - ApiRouterHandler: The router handler for chaining further route configurations.
	Delete(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T]

	// Patch registers a route handler specifically for HTTP PATCH requests with optional role-based access control.
	//
	// Parameters:
	// - handler: The handler function to process requests.
	// - path: The URL route path.
	// - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	// - ApiRouterHandler: The router handler for chaining further route configurations.
	Patch(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T]

	// Options registers a route handler specifically for HTTP OPTIONS requests with optional role-based access control.
	//
	// Parameters:
	// - handler: The handler function to process requests.
	// - path: The URL route path.
	// - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	// - ApiRouterHandler: The router handler for chaining further route configurations.
	Options(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T]

	// Head registers a route handler specifically for HTTP HEAD requests with optional role-based access control.
	//
	// Parameters:
	// - handler: The handler function to process requests.
	// - path: The URL route path.
	// - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	// - ApiRouterHandler: The router handler for chaining further route configurations.
	Head(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T]

	// Use adds a middleware function to the API router.
	// Middleware intercepts requests and can perform tasks like authentication, logging, etc.
	//
	// Parameters:
	// - middleware: The middleware function to apply to requests.
	// - name: The identifier for the middleware.
	//
	// Returns:
	// - ApiRouterHandler: The router handler for chaining further route configurations.
	Use(middleware ApiMiddleware[T], name string) ApiRouterHandler[T]

	// WithPrincipalService assigns a principal service to the router.
	// This service provides role-based access control and other principal-related features.
	//
	// Parameters:
	// - service: The principal service instance.
	//
	// Returns:
	// - ApiRouterHandler: The router handler for chaining further route configurations.
	WithPrincipalService(service *principal.PService[T]) ApiRouterHandler[T]

	WithErrorHandler(handler *ApiErrorHandler[T]) ApiRouterHandler[T]

	// StartServer starts the HTTP server with the configured routes and middleware.
	// This method blocks the current execution until the server terminates.
	StartServer()
}
