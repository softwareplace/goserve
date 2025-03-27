# goserve

`goserve` is a Go library designed to simplify the creation of backend applications or services that interact with
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

## Why Choose goserve?

Whether you're building lightweight services or complex applications, `goserve` provides a robust and
developer-friendly approach to creating HTTP servers. Its design ensures security, maintainability, and ease of use,
enabling you to focus on building features without worrying about boilerplate.

## Installation

Install `goserve` using:

```shell
go get -u github.com/softwareplace/goserve
```

## Getting Started

### Setting Up a Secure Backend Server as [example](./internal/main.go)

Configure a server with role-based access control, secure API authentication, and integrated middleware:

```go
package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/softwareplace/goserve/internal/handler"
	"github.com/softwareplace/goserve/internal/service/apiservice"
	"github.com/softwareplace/goserve/internal/service/login"
	"github.com/softwareplace/goserve/internal/service/provider"
	"github.com/softwareplace/goserve/logger"
	"github.com/softwareplace/goserve/security"
	"github.com/softwareplace/goserve/security/secret"
	"github.com/softwareplace/goserve/server"
	"os"
)

func init() {
	// Setup log system. Using nested-logrus-formatter -> https://github.com/antonfisher/nested-logrus-formatter?tab=readme-ov-file
	// Reload log file target reference based on `LOG_FILE_NAME_DATE_FORMAT`
	logger.LogSetup()
}

var (
	userPrincipalService = login.NewPrincipalService()
	securityService = security.New(
		"ue1pUOtCGaYS7Z1DLJ80nFtZ",
		userPrincipalService,
	)

	loginService = login.NewLoginService(securityService)
	secretProvider = provider.NewSecretProvider()

	secretService = secret.New(
		"./internal/secret/private.key",
		secretProvider,
		securityService,
	)
)

func main() {
	server.Default().
		LoginResourceEnabled(true).
		SecretKeyGeneratorResourceEnabled(true).
		LoginService(loginService).
		SecretService(secretService).
		SecurityService(securityService).
		EmbeddedServer(apiservice.Register).
		Get(apiservice.ReportCallerHandler, "/report/caller").
		SwaggerDocHandler("./internal/resource/pet-store.yaml").
		StartServer()
}

```

---

### Context Path Configuration

By default, the server runs at `http://localhost:8080/`. You can change the port and context path using the
following environment variables:

| Name                      | Required | Default    | Description                                                                                        |
|---------------------------|----------|------------|----------------------------------------------------------------------------------------------------|
| CONTEXT_PATH              | No       | /          | The base path for the application API.                                                             |
| PORT                      | No       | 8080       | The port on which the application runs.                                                            |
| B_CRYPT_COST              | No       | 10         | The cost factor for bcrypt hashing.                                                                |
| LOG_DIR                   | No       | ./.log     | The directory path where log files will be stored.                                                 |
| LOG_APP_NAME              | No       |            | The application name tha used to define the log file name.                                         |
| LOG_REPORT_CALLER         | No       | false      | [Logging Method Name](https://github.com/sirupsen/logrus?tab=readme-ov-file#logging-method-name)   |
| LOG_FILE_NAME_DATE_FORMAT | No       | 2006-01-02 | Defines the date format (YYYY-MM-DD) used for naming log files and tracking the current date       |
| JWT_ISSUER                | No       |            | Issuer returns the identifier of the entity responsible for issuing the JWT tokens in the service. |

### Advanced Configuration with Code Generation

Here's the updated text with the requested change:

---

### Configuration

Below is a required YAML configuration for customizing the code generation process to ensure compatibility with the
`goserve` library. This configuration defines the necessary templates and options for generating Go code that works
seamlessly with the library.

```yaml
package: gen

generate:
  gorilla-server: true
  models: true
output: ./gen/api.gen.go
output-options:
  user-templates:
    imports.tmpl: https://raw.githubusercontent.com/softwareplace/goserve/refs/heads/main/resource/templates/imports.tmpl
    param-types.tmpl: https://raw.githubusercontent.com/softwareplace/goserve/refs/heads/main/resource/templates/param-types.tmpl
    request-bodies.tmpl: https://raw.githubusercontent.com/softwareplace/goserve/refs/heads/main/resource/templates/request-bodies.tmpl
    typedef.tmpl: https://raw.githubusercontent.com/softwareplace/goserve/refs/heads/main/resource/templates/typedef.tmpl
    gorilla/gorilla-register.tmpl: https://raw.githubusercontent.com/softwareplace/goserve/refs/heads/main/resource/templates/gorilla/gorilla-register.tmpl
    gorilla/gorilla-middleware.tmpl: https://raw.githubusercontent.com/softwareplace/goserve/refs/heads/main/resource/templates/gorilla/gorilla-middleware.tmpl
    gorilla/gorilla-interface.tmpl: https://raw.githubusercontent.com/softwareplace/goserve/refs/heads/main/resource/templates/gorilla/gorilla-interface.tmpl
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

## Api Test

### Start server

```shell
go run internal/main.go
```

> Trying to access http://localhost:8080/swagger/index.html, it might respond with
>
> ![swagger-ui-example.png](internal/doc/swagger-ui-example.png)

### Start server in protected mode

```shell
PROTECTED_API=true go run internal/main.go
```

> Trying to access http://localhost:8080/swagger/index.html, it might respond with 401
>
> ```json
> {
>  "message": "You are not allowed to access this resource",
>  "statusCode": 401,
>  "timestamp": 1742781093916
> }
>```

> Try now passing jwt header of `X-Api-Key` and might respond with the swagger html entry.
> ````shell
> curl -X 'GET' \
> 'http://localhost:8080/swagger/index.html' \
> -H 'accept: application/json' \
> -H 'X-Api-Key: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcGlLZXkiOiJTL2VTYzVDQ3Jub1laaDAyU2pLdFVsSzFXdmRaaVA1OXpFUU9jNE54K0pjL1c1dkhMa0tndE1ueExHN3dKTUwvIiwiY2xpZW50IjoiU29mdHdhcmUgUGxhY2UiLCJleHAiOjMwMzM5MzczNTcsInNjb3BlIjpbIkhFMSs0cEVwM3YzZFBzWXNLa3FLMGkzdiswSjMvYjFVN01YQkx3ZzhxQ0E9IiwiR2lQWUVNU1IvK1BjNUdaTm9OcUpqZDRkS1FZbjZ6QzBMbmdYTHVxdFc4VzkiLCJjY294TWNaT0tEZ0srTUZuend0YWFEWXgxaEtPSVlKNDl3PT0iLCJxOWRHb3V5bTBxZWxvV1V4bElKZ2Y1U3l6UnIrU3YwWWwvVT0iLCJNOHdrRkN3cmZpeVBKc2hjb3NrQU5GS0RZZ2ZxRnJOWXkwVmljOEdlM3dPSyJdfQ.n5_8kp3nNqXOAZVB73GCIXcv61gNyyihqz6xDIjIA0k'
> ````

## Important Notes

- **Gorilla Mux Support**: The `gorilla-server` option ensures compatibility with `gorilla/mux`, one of the most popular
  Go routers.
- **Middleware Order**: The `apply-gorilla-middleware-first-to-last` option guarantees correct middleware application
  order.
- **Custom Templates**: Utilize hosted templates to align generated code with `goserve` conventions.

By following these steps, you can harness the full potential of `goserve` to build secure, scalable, and maintainable
HTTP servers in Go.

--- 
