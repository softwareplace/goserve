package main

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/security/principal"
	"github.com/softwareplace/http-utils/server"
)

func main() {
	contextBuilder := func(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) (doNext bool) {
		ctx.RequestData = &api_context.DefaultContext{}
		return true
	}

	var service principal.PService[*api_context.DefaultContext]
	service = &principalServiceImpl{}

	var errorHandler server.ApiErrorHandler[*api_context.DefaultContext]
	errorHandler = &errorHandlerImpl{}

	server.Default(contextBuilder).
		WithPrincipalService(&service).
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
