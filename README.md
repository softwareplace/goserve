# http-utils

> This library provides a simple and efficient way to create and start HTTP servers in Go. It abstracts away the
> boilerplate code required to set up a server, so you can focus on implementing your application's logic. Whether you
> are
> building a small service or a complex application, `http-utils` makes it easy to get started quickly while following
> best practices.
>
> This is a flexible and extensible API router structure for handling different HTTP methods (GET, POST, PUT, DELETE,
> etc.) and offers middleware, authentication, role-based access control, and error handling features. Here's a
> breakdown
> of the core components:

````shell
go get -u github.com/softwareplace/http-utils
````

```go
package main

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/error_handler"
	"github.com/softwareplace/http-utils/example/gen"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/security/principal"
	"github.com/softwareplace/http-utils/server"
	"log"
	"os"
	"time"
)

type loginServiceImpl struct {
	securityService security.ApiSecurityService[*api_context.DefaultContext]
}

func (l *loginServiceImpl) SecurityService() security.ApiSecurityService[*api_context.DefaultContext] {
	return l.securityService
}

func (l *loginServiceImpl) Login(user server.LoginEntryData) (*api_context.DefaultContext, error) {
	result := &api_context.DefaultContext{}
	result.SetRoles("api:example:user", "api:example:admin")
	return result, nil
}

func (l *loginServiceImpl) TokenDuration() time.Duration {
	return time.Minute * 15
}

type secretProviderImpl []struct{}

func (s *secretProviderImpl) Get(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) (string, error) {
	return "", nil
}

type principalServiceImpl struct {
}

func (d *principalServiceImpl) LoadPrincipal(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) bool {
	if ctx.Authorization == "" {
		return false

	}

	context := api_context.NewDefaultCtx()
	ctx.Principal = &context
	return true
}

type errorHandlerImpl struct {
}

func (p *errorHandlerImpl) Handler(ctx *api_context.ApiRequestContext[*api_context.DefaultContext], _ error, source string) {
	if source == server.ErrorHandlerWrapper {
		ctx.InternalServerError("Internal server error")
	}

	if source == server.SecurityValidatorResourceAccess {
		ctx.Unauthorized()
	}
}

type _service struct {
}

func (s *_service) PostLoginRequest(body gen.LoginRequest, ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {

}

func (s *_service) GetTestRequest(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	message := "It's working"
	code := 200
	success := true
	timestamp := 1625867200

	response := gen.BaseResponse{
		Message:   &message,
		Code:      &code,
		Success:   &success,
		Timestamp: &timestamp,
	}

	ctx.Response(response, 200)
}

func (s *_service) GetTestVersionRequest(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	message := "Test v2 it's working"
	code := 200
	success := true
	timestamp := 1625867200

	response := gen.BaseResponse{
		Message:   &message,
		Code:      &code,
		Success:   &success,
		Timestamp: &timestamp,
	}

	ctx.Response(response, 200)
}

func (s *_service) PostTestVersionRequest(body gen.PostTestRequest, ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	timestamp := 1625867200
	message := "Provided message was => " + body.Message
	response := gen.BaseResponse{
		Message:   &message,
		Code:      &body.Code,
		Success:   &body.Success,
		Timestamp: &timestamp,
	}

	ctx.Response(response, 200)
}

func (s *_service) GetTestV2(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	message := "Test v2 it's working"
	code := 200
	success := true
	timestamp := 1625867200

	response := gen.BaseResponse{
		Message:   &message,
		Code:      &code,
		Success:   &success,
		Timestamp: &timestamp,
	}

	ctx.Response(response, 200)
}

func main() {

	var userPrincipalService principal.PService[*api_context.DefaultContext]
	userPrincipalService = &principalServiceImpl{}

	var errorHandler error_handler.ApiErrorHandler[*api_context.DefaultContext]
	errorHandler = &errorHandlerImpl{}

	securityService := security.ApiSecurityServiceBuild(
		"ue1pUOtCGaYS7Z1DLJ80nFtZ",
		userPrincipalService,
	)

	secretProvider := &secretProviderImpl{}

	secretHandler := security.ApiSecretAccessHandlerBuild(
		"./example/secret/private.key",
		secretProvider,
		securityService,
	)

	secretHandler.DisableForPublicPath(true)

	for _, arg := range os.Args {
		if arg == "--d" || arg == "-d" {
			log.Println("Setting public path requires access with api secret key.")
			secretHandler.DisableForPublicPath(false)
		}
	}

	loginService := &loginServiceImpl{
		securityService: securityService,
	}

	server.Default().
		WithLoginResource(loginService).
		EmbeddedServer(gen.ApiResourceHandler(&_service{})).
		SwaggerDocHandler("example/resource/swagger.yaml").
		//SwaggerDocHandler("example/resource/swagger-with-api-key.yaml").
		RegisterMiddleware(secretHandler.HandlerSecretAccess, security.ApiSecretAccessHandlerName).
		RegisterMiddleware(securityService.AuthorizationHandler, security.ApiSecurityHandlerName).
		WithErrorHandler(errorHandler).
		NotFoundHandler().
		StartServer()
}

```

### Types

1. **ApiContextHandler**: A function type that handles requests within an API context, which includes the principal
   context for managing user roles and authentication.

2. **ApiMiddleware**: A function type that is used for middleware, which can intercept and modify requests before they
   reach the handler.

3. **RouterHandler**: A function type for handling route registrations with the `mux.Router`.

4. **ApiRouterHandler**: An interface that outlines the API router's capabilities, including methods for route handling,
   middleware registration, authentication, and documentation.

### Key Methods

- **PublicRouter**: Registers a public route, meaning it doesn't require authentication.

- **Add**: Registers a route with optional role-based access control.

- **Get**, **Post**, **Put**, **Delete**, **Patch**, **Options**, **Head**: Methods for registering routes with
  corresponding HTTP methods and optional role-based access control.

- **RegisterMiddleware**: Adds a middleware function to the API router.

- **RegisterCustomMiddleware**: Allows registering a custom middleware for additional flexibility in handling requests.

- **WithPrincipalService**: Associates a principal service with the router to handle role-based access control.

- **WithErrorHandler**: Registers a custom error handler to manage API errors.

- **WithLoginResource**: Configures a login resource for the router to handle user authentication.

- **Router**: Provides access to the underlying `mux.Router` instance, allowing more advanced routing or custom routes.

- **SwaggerDocProvider**: Configures the Swagger UI and serves the OpenAPI 3.0 documentation for the API.

- **SwaggerDocHandler**: Loads Swagger documentation from a file and serves it via the API.

- **NotFoundHandler**: Configures a handler for undefined routes, useful for providing custom 404 behavior.

### Example Use Case

This structure is ideal for creating APIs that need fine-grained control over routing, middleware, authentication, and
documentation. It allows the creation of public and private routes, integration with authentication and role management
systems, and the serving of API documentation via Swagger.

The router methods allow chaining, making it easy to configure routes and their associated behaviors. The flexible
middleware registration allows for a highly customizable request-handling pipeline. Additionally, integrating Swagger
documentation helps make the API self-descriptive and easier to work with for developers consuming the API.

Got it! Here's a documentation overview for the library, explaining its purpose and usage.

---

## Swagger Documentation Integration

This package provides functionality to integrate Swagger API documentation into a Go server, allowing you to dynamically
load and serve Swagger files for your API. It leverages the **`kin-openapi`** library to parse Swagger files and *
*`http-swagger`** to serve the interactive Swagger UI.

### Features:

- **Swagger File Parsing**: Load and parse Swagger specification files (both JSON and YAML formats).
- **Dynamic Path Adjustment**: Prepend custom context paths to all API paths in the Swagger spec.
- **Interactive UI**: Serve Swagger's interactive documentation UI, where users can explore and test API endpoints.
- **JSON Endpoint**: Expose the Swagger documentation as a raw JSON file at a configurable endpoint.

### Key Components:

1. **SwaggerDocLoader**: This function loads a Swagger file from the filesystem and parses it into an OpenAPI object.
   The resulting object can be used to serve the API documentation or modify the Swagger spec as needed.

2. **SwaggerDocProvider**: This method integrates Swagger documentation into the Go server’s routing system. It allows
   users to specify a custom function that loads Swagger specs and configures the server to serve the interactive UI and
   raw JSON file. It also adjusts paths in the Swagger file to be prefixed with a custom `ContextPath`.

3. **SwaggerDocHandler**: A convenience function that loads a Swagger file and serves it through the API router. It
   combines the process of loading the Swagger file and setting up necessary routes into one step.

4. **handleSwaggerJSON**: This function serves the Swagger JSON file directly as an HTTP response. It’s a helper for
   serving the raw API documentation in JSON format.

### Configuration:

- **`ContextPath`**: The path prefix that will be applied to all API routes in the Swagger documentation. For example,
  if `ContextPath` is set to `/api/`, all Swagger paths will be prefixed with `/api/`.
- **Swagger File**: The Swagger specification file (usually `swagger.json` or `swagger.yaml`) that describes the API's
  endpoints, parameters, and responses.

### Usage:

1. **Load Swagger Documentation**: You can use the `SwaggerDocLoader` function to load a Swagger spec from a file. This
   returns an OpenAPI object that can be further manipulated if needed.

2. **Serve Swagger UI**: Once the Swagger file is loaded, you can call the `SwaggerDocProvider` method to integrate it
   with your server. This will automatically configure routes to serve both the Swagger UI and the raw Swagger JSON.

3. **Custom Paths**: By default, the paths in the Swagger documentation will be adjusted to include the `ContextPath`.
   For instance, `/resource` will become `/api/resource` if `ContextPath` is set to `/api/`.

4. **Interactive UI**: The Swagger UI allows users to interact with the API documentation, test endpoints, and view
   detailed information about each API operation.
   Here's an improved version of the documentation with clearer instructions and additional context:

---

## Code Generation Compatibility with `http-utils` Library

To ensure that the generated code is compatible with this library, you need to configure and generate your API code
using the [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) tool. This tool converts OpenAPI specifications
into Go code that works seamlessly with the `http-utils` library.

### Configuration

Below is an example of how to configure the code generation process using a custom YAML config file. This configuration
defines the necessary templates and options for generating Go code that works with the `http-utils` library.

```yaml
package: gen

generate:
  gorilla-server: true  # Generate server code using Gorilla Mux
  models: true           # Generate models from the Swagger spec
#  embedded-spec: true   # Optionally embed the Swagger spec directly in the generated code
output: ./gen/api.gen.go  # Output path for the generated code

output-options:
  user-templates:
    param-types.tmpl: https://raw.githubusercontent.com/softwareplace/http-utils/refs/heads/main/resource/templates/param-types.tmpl
    additional-properties.tmpl: https://raw.githubusercontent.com/softwareplace/http-utils/refs/heads/main/resource/templates/additional-properties.tmpl
    request-bodies.tmpl: https://raw.githubusercontent.com/softwareplace/http-utils/refs/heads/main/resource/templates/request-bodies.tmpl
    typedef.tmpl: https://raw.githubusercontent.com/softwareplace/http-utils/refs/heads/main/resource/templates/typedef.tmpl
    gorilla/gorilla-register.tmpl: https://raw.githubusercontent.com/softwareplace/http-utils/refs/heads/main/resource/templates/gorilla/gorilla-register.tmpl
    gorilla/gorilla-middleware.tmpl: https://raw.githubusercontent.com/softwareplace/http-utils/refs/heads/main/resource/templates/gorilla/gorilla-middleware.tmpl
    gorilla/gorilla-interface.tmpl: https://raw.githubusercontent.com/softwareplace/http-utils/refs/heads/main/resource/templates/gorilla/gorilla-interface.tmpl

compatibility:
  apply-gorilla-middleware-first-to-last: true  # Ensure middleware is applied in the correct order
```

### Steps for Code Generation

1. **Create the Configuration File**:
   Create a YAML file based on the structure above. Make sure to define paths for the required templates and set the
   appropriate options for code generation.

2. **Run oapi-codegen**:
   Once you have your configuration file ready, use the following command to generate the Go code:

   ```shell
   oapi-codegen --config path/to/config.yaml path/to/swagger.yaml
   ```

    - Replace `path/to/config.yaml` with the path to your config file.
    - Replace `path/to/swagger.yaml` with the path to your OpenAPI/Swagger specification file.

3. **Generated Code**:
   The generated Go code will be output to the specified `output` path (e.g., `./gen/api.gen.go`), ready to be
   integrated with your application.

### Important Notes:

- **Template Customization**: This configuration specifies custom templates that align the generated code with the
  structure expected by the `http-utils` library. The templates are hosted in a GitHub repository and are referenced by
  URLs.

- **Gorilla Mux Server**: The `gorilla-server` option ensures that the generated server code is compatible with the
  Gorilla Mux router, which is supported by `http-utils`.

- **Middleware**: The `apply-gorilla-middleware-first-to-last` option ensures that the middleware is applied in the
  correct order to all generated routes.

By following these steps and using the correct configuration, you can generate Go code that works seamlessly with the
`http-utils` library and the associated templates.

---

This improved documentation provides more context on how to configure and generate the code, highlighting important
aspects such as templates and the middleware application order.

### Example Workflow:

1. Load a Swagger file using `SwaggerDocLoader`.
2. Integrate the Swagger documentation into the server using `SwaggerDocProvider`, providing a custom path for API
   routes (via `ContextPath`).
3. Serve the Swagger UI and JSON file at specified endpoints.

### Routes:

- **Swagger UI**: Accessible at `{ContextPath}/swagger/`.
- **Swagger JSON**: Available at `{ContextPath}/doc.json`.

---

> This package simplifies the process of exposing API documentation in a Go server, with support for custom paths and
> easy integration with the Swagger UI. It's designed to be flexible and easy to use, enabling developers to document
> and share their API endpoints quickly.


> You can also follow the [example](./example/main.go)


