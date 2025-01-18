package main

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/error_handler"
	"github.com/softwareplace/http-utils/example/impl"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/security/principal"
	"github.com/softwareplace/http-utils/server"
	"os"
	"time"
)

func main() {
	var service principal.PService[*api_context.DefaultContext]
	service = &impl.PrincipalServiceImpl{}

	var errorHandler error_handler.ApiErrorHandler[*api_context.DefaultContext]
	errorHandler = &impl.ErrorHandlerImpl{}

	securityService := security.ApiSecurityServiceBuild(
		"ue1pUOtCGaYS7Z1DLJ80nFtZ",
		&service,
	)

	loader := func(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) (string, error) {
		return "", nil
	}

	secretHandler := security.ApiSecretAccessHandlerBuild(
		"./example/secret/private.key",
		loader,
		securityService,
	)

	secretHandler.DisableForPublicPath(true)

	loginService := impl.New(securityService)

	server.Default().
		RegisterMiddleware(secretHandler.HandlerSecretAccess, security.ApiSecretAccessHandlerName).
		RegisterMiddleware(securityService.AuthorizationHandler, security.ApiSecurityHandlerName).
		WithLoginResource(&loginService).
		PublicRouter(isWorking, "example", "GET").
		PublicRouter(shutdown, "shutdown", "GET").
		Get(isWorkingV2, "example/v2", "GET", "example:v2").
		WithErrorHandler(&errorHandler).
		StartServer()
}

func isWorkingV2(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Response(map[string]string{"message": "It's working"}, 200)
}

func isWorking(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Response(map[string]string{"message": "It's working"}, 200)
}

func shutdown(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Response(map[string]string{"message": "Shutting down in one second"}, 200)

	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}
