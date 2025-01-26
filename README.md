Here’s an improved version of your documentation with explicit references to the `gorilla/mux` library:

---

# http-utils

`http-utils` is a Go library designed to simplify the creation of backend applications or services that interact with
HTTP requests. It leverages the powerful `gorilla/mux` router to provide flexibility, performance, and scalability while
adhering to best practices in server development.

## Key Features

- **Backend Application Server**: Quickly set up a backend server with support for security, role-based access control,
  and efficient resource handling.
- **Enhanced Security**: Protect your application with an `apiSecret` mechanism, powered by `private.key` and
  `public.key` for authentication and validation.
- **Swagger-UI Integration**: Simplify API documentation setup
  using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) and streamline request data handling (e.g., body,
  query parameters, headers).
- **Router Flexibility**: Seamlessly handle HTTP methods (GET, POST, PUT, DELETE) with `gorilla/mux`, enabling powerful
  URL matching and routing.
- **Built-in Middleware**: Apply middleware for authentication, role-based access control, and error handling with
  minimal setup.

## Why Choose http-utils?

Whether you're building lightweight services or complex applications, `http-utils` provides a robust and
developer-friendly approach to creating HTTP servers. Its design ensures security, maintainability, and ease of use,
enabling you to focus on building features without worrying about boilerplate.

## Installation

Install `http-utils` using:

```shell
go get -u github.com/softwareplace/http-utils
```

## Getting Started

### Setting Up a Secure Backend Server

Configure a server with role-based access control, secure API authentication, and integrated middleware:

- Full example [here](./example/main.go)

```go
package main

import (
	"github.com/softwareplace/http-utils/example/gen"
	"github.com/softwareplace/http-utils/server"
)

func main() {
	server.Default().
		LoginResourceEnabled(true). // Enable login resource
		ApiSecretKeyGeneratorResourceEnabled(true). // Enable API secret key generator resource
		WithLoginResource(loginService). // Attach the login service
		WithApiKeyGeneratorResource(loginService). // Attach API key generator service
		EmbeddedServer(gen.ApiResourceHandler(&_service{})). // Add embedded API resource handler
		SwaggerDocHandler("example/resource/pet-store.yaml"). // Serve Swagger-UI
		WithApiSecretAccessHandler(secretHandler). // Configure secret access handler
		WithApiSecurityService(securityService). // Set up security service
		WithErrorHandler(errorHandler). // Define custom error handler
		NotFoundHandler(). // Handle 404 errors
		StartServer() // Start the server
}
```

---

### Context Path Configuration

By default, the server runs at `http://localhost:8080/api/app/v1/`. You can change the port and context path using the
following environment variables:

| Name         | Required | Default      |
|--------------|----------|--------------|
| CONTEXT_PATH | No       | /api/app/v1/ |
| PORT         | No       | 8080         |

### Advanced Configuration with Code Generation

Here's the updated text with the requested change:

---

### Configuration

Below is a required YAML configuration for customizing the code generation process to ensure compatibility with the
`http-utils` library. This configuration defines the necessary templates and options for generating Go code that works
seamlessly with the library.

```yaml
package: gen

generate:
  gorilla-server: true
  models: true
output: ./gen/api.gen.go
output-options:
  user-templates:
    imports.tmpl: https://raw.githubusercontent.com/softwareplace/http-utils/refs/heads/main/resource/templates/imports.tmpl
    param-types.tmpl: https://raw.githubusercontent.com/softwareplace/http-utils/refs/heads/main/resource/templates/param-types.tmpl
    request-bodies.tmpl: https://raw.githubusercontent.com/softwareplace/http-utils/refs/heads/main/resource/templates/request-bodies.tmpl
    typedef.tmpl: https://raw.githubusercontent.com/softwareplace/http-utils/refs/heads/main/resource/templates/typedef.tmpl
    gorilla/gorilla-register.tmpl: https://raw.githubusercontent.com/softwareplace/http-utils/refs/heads/main/resource/templates/gorilla/gorilla-register.tmpl
    gorilla/gorilla-middleware.tmpl: https://raw.githubusercontent.com/softwareplace/http-utils/refs/heads/main/resource/templates/gorilla/gorilla-middleware.tmpl
    gorilla/gorilla-interface.tmpl: https://raw.githubusercontent.com/softwareplace/http-utils/refs/heads/main/resource/templates/gorilla/gorilla-interface.tmpl
```

---

### Steps for Code Generation

1. **Create the YAML Configuration File**: Use the structure above to define paths for templates and set options for
   `oapi-codegen`.

2. **Run Code Generation**: Generate the Go code with:

   ```shell
   oapi-codegen --config path/to/config.yaml path/to/swagger.yaml
   ```

   Replace `path/to/config.yaml` with your configuration file path and `path/to/swagger.yaml` with your OpenAPI/Swagger
   spec.

3. **Integrate the Generated Code**: The output (e.g., `./gen/api.gen.go`) will be ready for use in your application.

---

## Important Notes

- **Gorilla Mux Support**: The `gorilla-server` option ensures compatibility with `gorilla/mux`, one of the most popular
  Go routers.
- **Middleware Order**: The `apply-gorilla-middleware-first-to-last` option guarantees correct middleware application
  order.
- **Custom Templates**: Utilize hosted templates to align generated code with `http-utils` conventions.

By following these steps, you can harness the full potential of `http-utils` to build secure, scalable, and maintainable
HTTP servers in Go.

--- 
