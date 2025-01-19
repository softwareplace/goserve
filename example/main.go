package main

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/error_handler"
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

func main() {

	var service principal.PService[*api_context.DefaultContext]
	service = &principalServiceImpl{}

	var errorHandler error_handler.ApiErrorHandler[*api_context.DefaultContext]
	errorHandler = &errorHandlerImpl{}

	securityService := security.ApiSecurityServiceBuild(
		"ue1pUOtCGaYS7Z1DLJ80nFtZ",
		service,
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

	go func() {
		log.Println("Application will be shut down in 5 seconds.")
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()

	server.Default().
		RegisterMiddleware(secretHandler.HandlerSecretAccess, security.ApiSecretAccessHandlerName).
		RegisterMiddleware(securityService.AuthorizationHandler, security.ApiSecurityHandlerName).
		WithLoginResource(loginService).
		PublicRouter(isWorking, "test", "GET").
		Get(isWorkingV2, "test/v2", "api:example:admin").
		WithErrorHandler(errorHandler).
		StartServer()
}

func isWorkingV2(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Response(map[string]string{"message": "It's working"}, 200)
}

func isWorking(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Response(map[string]string{"message": "It's working"}, 200)
}
