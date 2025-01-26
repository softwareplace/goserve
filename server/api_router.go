package server

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/error_handler"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/security/principal"
	"net/http"
)

type ApiContextHandler[T api_context.ApiPrincipalContext] func(ctx *api_context.ApiRequestContext[T])

type ApiMiddleware[T api_context.ApiPrincipalContext] func(*api_context.ApiRequestContext[T]) (doNext bool)

type RouterHandler func(*mux.Router)

type ApiRouterHandler[T api_context.ApiPrincipalContext] interface {
	// PublicRouter registers a public route handler that does not require authentication or authorization.
	// It allows unrestricted access from any client.
	//
	// Parameters:
	//   - handler: The handler function to process requests.
	//   - path: The URL route path.
	//   - method: The HTTP method for the route.
	//
	// Returns:
	//   - ApiRouterHandler: The router handler for chaining further route configurations.
	PublicRouter(handler ApiContextHandler[T], path string, method string) ApiRouterHandler[T]

	// Add registers a route handler with optional role-based access control.
	// This method is used to define routes and assign roles required to access them.
	//
	// Parameters:
	//   - handler: The handler function to process requests.
	//   - path: The URL route path.
	//   - method: The HTTP method for the route.
	//   - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	//   - ApiRouterHandler: The router handler for chaining further route configurations.
	Add(handler ApiContextHandler[T], path string, method string, requiredRoles ...string) ApiRouterHandler[T]

	// Get registers a route handler specifically for HTTP GET requests with optional role-based access control.
	//
	// Parameters:
	//   - handler: The handler function to process requests.
	//   - path: The URL route path.
	//   - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	//   - ApiRouterHandler: The router handler for chaining further route configurations.
	Get(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T]

	// Post registers a route handler specifically for HTTP POST requests with optional role-based access control.
	//
	// Parameters:
	//   - handler: The handler function to process requests.
	//   - path: The URL route path.
	//   - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	//   - ApiRouterHandler: The router handler for chaining further route configurations.
	Post(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T]

	// Put registers a route handler specifically for HTTP PUT requests with optional role-based access control.
	//
	// Parameters:
	//   - handler: The handler function to process requests.
	//   - path: The URL route path.
	//   - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	//   - ApiRouterHandler: The router handler for chaining further route configurations.
	Put(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T]

	// Delete registers a route handler specifically for HTTP DELETE requests with optional role-based access control.
	//
	// Parameters:
	//   - handler: The handler function to process requests.
	//   - path: The URL route path.
	//   - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	//   - ApiRouterHandler: The router handler for chaining further route configurations.
	Delete(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T]

	// Patch registers a route handler specifically for HTTP PATCH requests with optional role-based access control.
	//
	// Parameters:
	//   - handler: The handler function to process requests.
	//   - path: The URL route path.
	//   - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	//   - ApiRouterHandler: The router handler for chaining further route configurations.
	Patch(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T]

	// Options registers a route handler specifically for HTTP OPTIONS requests with optional role-based access control.
	//
	// Parameters:
	//   - handler: The handler function to process requests.
	//   - path: The URL route path.
	//   - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	//   - ApiRouterHandler: The router handler for chaining further route configurations.
	Options(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T]

	// Head registers a route handler specifically for HTTP HEAD requests with optional role-based access control.
	//
	// Parameters:
	//   - handler: The handler function to process requests.
	//   - path: The URL route path.
	//   - requiredRoles: A list of roles required to access the route (optional).
	//
	// Returns:
	//   - ApiRouterHandler: The router handler for chaining further route configurations.
	Head(handler ApiContextHandler[T], path string, requiredRoles ...string) ApiRouterHandler[T]

	// RegisterMiddleware adds a middleware function to the API router.
	// Middleware intercepts requests and can perform tasks like authentication, logging, etc.
	//
	// Parameters:
	//   - middleware: The middleware function to apply to requests.
	//   - name: The identifier for the middleware.
	//
	// Returns:
	//   - ApiRouterHandler: The router handler for chaining further route configurations.
	RegisterMiddleware(middleware ApiMiddleware[T], name string) ApiRouterHandler[T]

	// WithApiSecretAccessHandler assigns an API secret access handler to the router.
	// This middleware provides an additional layer of security by validating API secret keys
	// on incoming requests before they are processed by specific route handlers.
	//
	// Parameters:
	//   - apiSecretAccessHandler: An implementation of ApiSecretAccessHandler[T] responsible for the logic
	//				 to validate and handle API secret access.
	//
	// Returns:
	//   - ApiRouterHandler: The router handler for chaining further route configurations.
	WithApiSecretAccessHandler(apiSecretAccessHandler security.ApiSecretAccessHandler[T]) ApiRouterHandler[T]

	// WithApiSecurityService assigns a security service to the API router.
	// This service is responsible for handling various security aspects such as authentication,
	// authorization, and token validation for incoming requests. By assigning a security service,
	// the API router ensures that all operations comply with the defined security policies.
	//
	// Parameters:
	//   - apiSecurityService: An instance of ApiSecurityService[T] that implements the security logic.
	//
	// Returns:
	//   - ApiRouterHandler: The router handler for chaining additional route or service configurations.
	WithApiSecurityService(apiSecurityService security.ApiSecurityService[T]) ApiRouterHandler[T]

	// RegisterCustomMiddleware is a method that allows for registering a custom middleware function
	// with the API router. This function wraps an HTTP handler and can be used to implement custom
	// functionality, such as modifying the request or response, logging, or adding additional
	// security checks.
	//
	// Parameters:
	//   - next: The next HTTP handler in the middleware chain. It represents the subsequent middleware
	//		 or the final handler that should be invoked after the custom middleware logic.
	//
	// Example usage:
	// ```go
	// router := NewApiRouter()
	// router.RegisterCustomMiddleware(func(next http.Handler) http.Handler {
	//	 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		 // Custom middleware logic here
	//		 log.Println("Before handler")
	//		 next.ServeHTTP(w, r) // Call the next handler
	//		 log.Println("After handler")
	//	 })
	// })
	// ```
	//
	// This method does not return any value and directly modifies the middleware chain for the API router.
	RegisterCustomMiddleware(func(next http.Handler) http.Handler) ApiRouterHandler[T]

	// WithPrincipalService assigns a principal service to the router.
	// This service provides role-based access control and other principal-related features.
	//
	// Parameters:
	//   - service: The principal service instance.
	//
	// Returns:
	//   - ApiRouterHandler: The router handler for chaining further route configurations.
	WithPrincipalService(service principal.PService[T]) ApiRouterHandler[T]

	// WithErrorHandler assigns a custom error handler to the router.
	// This handler is used to process API errors and provide appropriate responses.
	//
	// Parameters:
	//   - handler: An instance of ApiErrorHandler that defines custom error-handling logic.
	//
	// Returns:
	//   - ApiRouterHandler: The router handler for chaining further route configurations.
	WithErrorHandler(handler error_handler.ApiErrorHandler[T]) ApiRouterHandler[T]

	// WithLoginResource sets the login service for the API router and registers a public
	// route for the login functionality. This route is accessible without requiring any
	// authentication or additional middleware.
	//
	// Parameters:
	//   - loginService: A pointer to the LoginService instance used to handle login functionality.
	//
	// This method internally registers a public POST endpoint at the path "/login" by linking
	// it to the `Login` handler method. The login service defined here allows managing and
	// validating user login flows.
	//
	// Returns:
	//   ApiRouterHandler[T]: This allows chaining additional configuration or service registrations.
	WithLoginResource(loginService LoginService[T]) ApiRouterHandler[T]

	// WithApiKeyGeneratorResource configures the ApiRouterHandler with a provided ApiKeyGeneratorService.
	// It registers a POST endpoint for generating API keys at the route path "/api-keys/generate".
	// This method ensures that the endpoint follows consistent naming conventions and best practices.
	//
	// Parameters:
	//   - apiKeyGeneratorService: The implementation of ApiKeyGeneratorService[T] responsible for generating API keys.
	//
	// Returns:
	//   - ApiRouterHandler[T]: The current ApiRouterHandler instance to allow method chaining.
	WithApiKeyGeneratorResource(apiKeyGeneratorService ApiKeyGeneratorService[T]) ApiRouterHandler[T]

	// Router retrieves the underlying *mux.Router instance.
	// This method provides direct access to the Gorilla Mux router, allowing you to add custom
	// routes, middleware, or additional configurations that are not covered by the ApiRouterHandler methods.
	//
	// Returns:
	//   - *mux.Router: The Gorilla Mux router instance.
	//
	// Example usage:
	// ```go
	// routerHandler := apiRouter.Router()
	// routerHandler.HandleFunc("/custom-path", customHandlerFunction)
	// ```
	Router() *mux.Router

	// RouterHandler assigns a custom RouterHandler interface to the API router.
	// This can be used to provide advanced or application-specific routing logic,
	// allowing greater flexibility in handling requests.
	//
	// Parameters:
	//   - handler: The custom RouterHandler instance to use for routing logic.
	//
	// Returns:
	//   - ApiRouterHandler: The router handler for chaining further route configurations.
	//
	// Example usage:
	// ```go
	// customHandler := NewCustomRouterHandler()
	// apiRouter.RouterHandler(customHandler)
	// ```
	//
	// This method provides a mechanism to integrate a custom router handler
	// implementation into the existing API router configuration pipeline.
	RouterHandler(handler RouterHandler) ApiRouterHandler[T]

	// EmbeddedServer allows embedding of an HTTP server within the application and enables
	// configuration of API routes, middleware, and other server settings through a handler function.
	//
	// Parameters:
	//   - handler: A function that accepts an ApiRouterHandler[T] and is used to configure the
	//			  server's routes, middleware, and other features. This provides a flexible way
	//			  to set up the server before starting it.
	//
	// Returns:
	//   - ApiRouterHandler[T]: The configured API router handler, allowing further chaining of
	//						configurations after embedding the server.
	//
	// Example usage:
	// ```go
	// apiRouter := api.EmbeddedServer(func(router api.ApiRouterHandler[T]) {
	//	 router.Get(myHandlerFunc, "/example-path")
	//	 router.Post(anotherHandlerFunc, "/post-path", "admin")
	// })
	// apiRouter.StartServer()
	// ```
	//
	// This can be used to create self-contained server setups for microservices, testing, or
	// other embedded server use cases.
	EmbeddedServer(handler func(ApiRouterHandler[T])) ApiRouterHandler[T]

	// SwaggerDocProvider sets up the Swagger UI and handles serving the
	// Swagger JSON documentation. It accepts a function to retrieve
	// the OpenAPI 3.0 specification.
	//
	// Parameters:
	//   - getSwagger: A function that returns the Swagger specification
	//	 (*openapi3.T) or an error if the specification cannot be loaded.
	//
	// Behavior:
	//   - The function retrieves the Swagger specification using getSwagger.
	//   - If retrieving the specification fails, it writes the error to
	//	 standard error and terminates the program.
	//   - Removes any server information from the Swagger specification to
	//	 prevent exposing unnecessary details.
	//   - Configures an HTTP handler to serve Swagger UI documentation.
	//   - Registers the Swagger UI handler under the "swagger/" path and
	//	 provides an endpoint for serving the raw Swagger JSON located at
	//	 "doc" using the context path configured.
	// Example:
	//   - Using generated swagger SwaggerDocProvider(gen.GetSwagger)
	SwaggerDocProvider(getSwagger func() (swagger *openapi3.T, err error)) ApiRouterHandler[T]

	// SwaggerDocHandler loads the Swagger documentation from a specified file into the API router.
	// This method reads the content of the provided file and serves it as the Swagger JSON.
	// It is useful for integrating pre-generated OpenAPI specifications into the API.
	//
	// Parameters:
	//   - swaggerFile: The file path to the Swagger/OpenAPI YAML file.
	//
	// Behavior:
	//   - Loads the Swagger specification from the specified file.
	//   - If the file cannot be read, the method returns an appropriate error or logs it,
	//	 depending on the implementation details.
	//   - Configures the API router to serve the Swagger definition via an HTTP handler.
	//
	// Returns:
	//   - ApiRouterHandler[T]: The router handler for chaining further configurations.
	//
	// Example usage:
	// ```go
	// router := NewApiRouter()
	// router.SwaggerDocHandler("path/to/swagger.yaml")
	// ```
	SwaggerDocHandler(swaggerFile string) ApiRouterHandler[T]

	// NotFoundHandler sets a custom handler for requests to undefined routes.
	// This method can be used to provide a user-friendly response or logging
	// for routes that are not registered within the API router.
	//
	// Behavior:
	//   - If SetupSwagger was invoked and the router matches with ContextPath, redirects to ContextPath+"swagger/index.html".
	//   - Registers a default HTTP handler for undefined routes.
	//   - Allows customization of the response for 404 Not Found errors.
	//   - If this method is not used, a standard "404 Page Not Found" is returned.
	//
	// Example usage:
	// ```go
	//    router := NewApiRouter()
	//    router.NotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	  	w.WriteHeader(http.StatusNotFound)
	//	    w.Write([]byte("Custom 404 Response"))
	//    }))
	// ```
	NotFoundHandler() ApiRouterHandler[T]

	// CustomNotFoundHandler sets a custom handler for handling 404 Not Found errors.
	// This method allows you to define your own response logic for requests to undefined routes.
	//
	// Parameters:
	//   - handler: A function that takes a http.ResponseWriter and *http.Request as arguments,
	//			and defines the logic for responding to unmatched routes.
	//
	// Returns:
	//   - ApiRouterHandler[T]: The router handler for chaining further route configurations.
	//
	// Example usage:
	// ```go
	// router := NewApiRouter()
	// router.CustomNotFoundHandler(func(w http.ResponseWriter, r *http.Request) {
	//	 w.WriteHeader(http.StatusNotFound)
	//	 w.Write([]byte("Custom 404 Page"))
	// })
	// ```
	CustomNotFoundHandler(handler func(w http.ResponseWriter, r *http.Request)) ApiRouterHandler[T]

	// StartServer starts the HTTP server with the configured routes and middleware.
	// This method blocks the current execution until the server terminates.
	StartServer()
}
