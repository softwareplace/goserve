# goserve

`goserve` is a Go library designed to simplify the creation of backend applications or services that interact with HTTP
requests. It leverages the powerful `gorilla/mux` router to provide flexibility, performance, and scalability while
adhering to best practices in server development.

---

## 🛠️ Prerequisites

Before using goserve, make sure you have the following installed:

1. **Go Programming Language**
   👉 [Install Go](https://go.dev/doc/install)
   ✅ Verify installation:
   ```bash
   go version
   ```

2. **Git** (Version Control System)
   👉 [Install Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
   ✅ Verify installation:
   ```bash
   git --version
   ```

3. **Environment Setup**
   Add Go binaries to your `PATH`:
   ```bash
   export PATH="$HOME/go/bin:$PATH"
   ```
   Add this line to your shell configuration file (e.g., `.bashrc`, `.zshrc`) to persist it.

4. **oapi-codegen Tool** (for OpenAPI/Swagger integration):
   ```bash
   go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.2.0
   ```

---

## 📦 Installation

Install the `goserve-generator` CLI tool:

```bash
go install github.com/softwareplace/goserve/cmd/goserve-generator@latest
```

Or add the library to your Go project:

```bash
go get -u github.com/softwareplace/goserve
```

---

## 🚀 Usage

Generate a new goserve project:

```bash
goserve-generator -n <project-name> -u <github-username> [-r true|false] [-gi true|false]
```

### Flags

| Flag      | Description                           | Required | Default |
| --------- | ------------------------------------- | -------- | ------- |
| `-n`      | Name of your project                  | ✅ Yes    |         |
| `-u`      | Your GitHub username                  | ✅ Yes    |         |
| `-r`      | Force replace existing files          | ❌ No     | false   |
| `-gi`     | Gi project initialization             | ❌ No     | true    |
| `-cgf`    | Template of the codegen config file   | ❌ No     |         |
| `-gsv`    | Use a specific version of goserver    | ❌ No     |         |
| `version` | Check current version                 | ❌ No     |         |
| `update`  | Update to the latest released version | ❌ No     |         |

### Example

```bash
goserve-generator -n goserve-example -u myuser -r true -gi false
```

---

## ✨ Key Features

- **Backend Application Server**: Kickstart a backend server with security, role-based access control, and scalable
  routing.
- **Enhanced Security**: Built-in support for API key authentication using `private.key` and `public.key`.
- **Swagger-UI Integration**: Built-in OpenAPI docs via `oapi-codegen`.
- **Router Flexibility**: Powered by `gorilla/mux` for clean, RESTful routing.
- **Built-in Middleware**: Support for authentication, role-checking, and structured error handling.

---

## 🛡️ Environment Variables

| Variable Name                   | Required? | Default      | Description                          |
| ------------------------------- | --------- | ------------ | ------------------------------------ |
| `CONTEXT_PATH`                  | No        | `/`          | Base path for all endpoints          |
| `PORT`                          | No        | `8080`       | Port the server listens on           |
| `API_SECRET_KEY`                | Yes*      |              | Used in encryption/authentication    |
| `API_PRIVATE_KEY`               | Yes*      |              | The private.key file path            |
| `B_CRYPT_COST`                  | No        | `10`         | Cost factor for bcrypt               |
| `LOG_DIR`                       | No        | `./.log`     | Where log files are stored           |
| `LOG_APP_NAME`                  | No        |              | Used in the log file naming          |
| `LOG_REPORT_CALLER`             | No        | `false`      | Enable method name reporting in logs |
| `LOG_FILE_NAME_DATE_FORMAT`     | No        | `2006-01-02` | Date format for log filenames        |
| `JWT_ISSUER`                    | No        |              | JWT issuer name                      |
| `JWT_CLAIMS_ENCRYPTION_ENABLED` | No        | `true`       | Encrypt claims inside JWT            |
| `SWAGGER_RESURCE_ENABLED`       | No        | `true`       | Enable swagger resource              |

\* Required only if using `security.Service`

---

## 🧪 Example: Secure Server Setup

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
)

func init() {
	logger.LogSetup()
}

var (
	userPrincipalService = login.NewPrincipalService()
	securityService      = security.New(
		userPrincipalService,
	)
	loginService   = login.NewLoginService(securityService)
	secretProvider = provider.NewSecretProvider()
	secretService  = secret.New(
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

## 🔧 Code Generation Config (oapi-codegen)

To customize code generation:

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

### Generate Code

```bash
oapi-codegen --config path/to/config.yaml path/to/swagger.yaml
```

---

## 📋 API Testing

### Start server

```bash
go run test/main.go
```

Open [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

### Start in protected mode

```bash
PROTECTED_API=true go run test/main.go
```

Without token:

```json
{
  "message": "You are not allowed to access this resource",
  "statusCode": 401,
  "timestamp": 1742781093916
}
```

With valid token:

```bash
curl -X GET 'http://localhost:8080/swagger/index.html' \
  -H 'accept: application/json' \
  -H 'X-Api-Key: <your-jwt-token>'
```

---

## 📚 Why Choose goserve?

Whether you're building microservices or full-stack applications, goserve provides a clean, secure, and production-ready
server foundation. With powerful integrations and easy configuration, goserve helps you focus on building features—not
boilerplate.

---

## 🧩 License

## 🎉 Welcome to the Open Source Community!

Open source is more than just code — it's a collaborative effort powered by you! Contributing to goserve means you're
joining a community of developers who believe in making software more accessible, secure, and scalable for everyone.

Thank you for being a part of this journey! Together, we create, innovate, and grow. 🚀
