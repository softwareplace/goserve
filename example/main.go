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

func (s *_service) PostLogin(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {

}

func (s *_service) GetTest(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
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

func embeddedHandler() func(handler server.ApiRouterHandler[*api_context.DefaultContext]) {
	return func(handler server.ApiRouterHandler[*api_context.DefaultContext]) {
		gen.ResourcesHandler(handler, &_service{})
	}
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
		EmbeddedServer(embeddedHandler()).
		RegisterMiddleware(secretHandler.HandlerSecretAccess, security.ApiSecretAccessHandlerName).
		RegisterMiddleware(securityService.AuthorizationHandler, security.ApiSecurityHandlerName).
		WithErrorHandler(errorHandler).
		StartServer()
}
