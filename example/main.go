package main

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/error_handler"
	"github.com/softwareplace/http-utils/example/impl"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/security/principal"
	"github.com/softwareplace/http-utils/server"
	"log"
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

	for _, arg := range os.Args {
		if arg == "--d" || arg == "-d" {
			log.Println("Setting public path requires access with api secret key.")
			secretHandler.DisableForPublicPath(false)
		}
	}

	loginService := impl.New(securityService)

	go func() {
		log.Println("Application will be shut down in 5 seconds.")
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()

	server.Default().
		RegisterMiddleware(secretHandler.HandlerSecretAccess, security.ApiSecretAccessHandlerName).
		RegisterMiddleware(securityService.AuthorizationHandler, security.ApiSecurityHandlerName).
		WithLoginResource(&loginService).
		PublicRouter(isWorking, "example", "GET").
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
