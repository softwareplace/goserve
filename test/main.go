package main

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/security/principal"
	"github.com/softwareplace/http-utils/server"
)

func main() {
	var service principal.PService[*api_context.DefaultContext]
	service = &principalServiceImpl{}

	var errorHandler server.ApiErrorHandler[*api_context.DefaultContext]
	errorHandler = &errorHandlerImpl{}

	securityService := security.ApiSecurityServiceBuild(
		"ue1pUOtCGaYS7Z1DLJ80nFtZ",
		&service,
	)

	loader := func(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) (string, error) {
		return "", nil
	}

	secretHandler := security.ApiSecretAccessHandlerBuild(
		"./test/secret/private.key",
		loader,
		securityService,
	)

	server.Default().
		RegisterMiddleware(secretHandler.Handler, security.ApiSecretAccessHandlerName).
		RegisterMiddleware(securityService.Handler, security.ApiSecurityHandlerName).
		PublicRouter(isWorking, "test", "GET").
		Get(isWorkingV2, "test/v2", "GET", "test:v2").
		WithErrorHandler(&errorHandler).
		StartServer()
}

func isWorkingV2(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Response(map[string]string{"message": "It's working"}, 200)
}

func isWorking(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Response(map[string]string{"message": "It's working"}, 200)
}
