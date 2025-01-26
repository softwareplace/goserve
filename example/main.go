package main

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/error_handler"
	"github.com/softwareplace/http-utils/example/gen"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/security/principal"
	"github.com/softwareplace/http-utils/server"
	"log"
	"net/http"
	"os"
	"time"
)

type loginServiceImpl struct {
	securityService security.ApiSecurityService[*api_context.DefaultContext]
}

func baseResponse(message string, status int) gen.BaseResponse {
	success := false
	timestamp := 1625867200

	response := gen.BaseResponse{
		Message:   &message,
		Code:      &status,
		Success:   &success,
		Timestamp: &timestamp,
	}
	return response
}

func (l *loginServiceImpl) SecurityService() security.ApiSecurityService[*api_context.DefaultContext] {
	return l.securityService
}

func (l *loginServiceImpl) GetApiJWTInfo(apiKeyEntryData server.ApiKeyEntryData,
	_ *api_context.ApiRequestContext[*api_context.DefaultContext],
) (security.ApiJWTInfo, error) {
	return security.ApiJWTInfo{
		Client:     apiKeyEntryData.ClientName,
		Key:        apiKeyEntryData.ClientId,
		Expiration: apiKeyEntryData.Expiration,
		Scopes: []string{
			"api:example:user",
			"api:example:admin",
			"read:pets",
			"write:pets",
			"api:key:generator",
		},
	}, nil
}

func (l *loginServiceImpl) OnGenerated(data security.JwtResponse,
	ctx api_context.SampleContext[*api_context.DefaultContext],
) {
	log.Printf("API KEY GENERATED: from %s - %v", ctx.AccessId, data)
}

func (l *loginServiceImpl) Login(user server.LoginEntryData) (*api_context.DefaultContext, error) {
	result := &api_context.DefaultContext{}
	result.SetRoles("api:example:user", "api:example:admin", "read:pets", "write:pets", "api:key:generator")
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

type _petService struct {
}

func (s *_petService) AddPetRequest(requestBody gen.Pet, ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Response(requestBody, http.StatusOK)
}

func (s *_petService) UpdatePetRequest(requestBody gen.Pet, ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Response(requestBody, http.StatusOK)
}

func (s *_petService) FindPetsByStatusRequest(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	message := "Pet not found"
	ctx.NotFount(baseResponse(message, http.StatusNotFound))
}

func (s *_petService) FindPetsByTagsRequest(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	message := "Pet not found"
	ctx.NotFount(baseResponse(message, http.StatusNotFound))
}

func (s *_petService) DeletePetRequest(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	message := "Pet not found"
	ctx.NotFount(baseResponse(message, http.StatusNotFound))
}

func (s *_petService) GetPetByIdRequest(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	message := "Pet not found"
	ctx.NotFount(baseResponse(message, http.StatusNotFound))
}

func (s *_petService) UpdatePetWithFormRequest(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	message := "Pet not found"
	ctx.NotFount(baseResponse(message, http.StatusNotFound))
}

type _userService struct {
}

func (s *_userService) PostLoginRequest(requestBody gen.LoginRequest, ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
}

func (s *_userService) CreateUserRequest(requestBody gen.User, ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Ok(requestBody)
}

func (s *_userService) CreateUsersWithListInputRequest(requestBody gen.CreateUsersWithListInputJSONBody, ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Ok(requestBody)
}

func (s *_userService) LogoutUserRequest(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Ok(baseResponse("Logout successful", http.StatusOK))
}

func (s *_userService) DeleteUserRequest(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Ok(baseResponse("User deleted", http.StatusOK))
}

func (s *_userService) GetUserByNameRequest(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.NotFount(baseResponse("User not found", http.StatusNotFound))
}

func (s *_userService) UpdateUserRequest(requestBody gen.User, ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Ok(requestBody)
}

type _fileService struct {
}

func (s *_fileService) UploadFileRequest(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.BadRequest("Failed to upload file")
}

type _inventoryService struct {
}

func (s *_inventoryService) GetInventoryRequest(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.NotFount(baseResponse("Inventory not found", http.StatusNotFound))
}

func (s *_inventoryService) PlaceOrderRequest(requestBody gen.Order, ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Response(requestBody, http.StatusOK)
}

func (s *_inventoryService) DeleteOrderRequest(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.Ok(baseResponse("Order deleted", http.StatusOK))
}

func (s *_inventoryService) GetOrderByIdRequest(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
	ctx.NotFount(baseResponse("Order not found", http.StatusNotFound))
}

type _service struct {
	_userService
	_petService
	_fileService
	_inventoryService
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
		WithApiKeyGeneratorResource(loginService, "api:key:generator").
		EmbeddedServer(gen.ApiResourceHandler(&_service{})).
		SwaggerDocHandler("example/resource/pet-store.yaml").
		RegisterMiddleware(secretHandler.HandlerSecretAccess, security.ApiSecretAccessHandlerName).
		RegisterMiddleware(securityService.AuthorizationHandler, security.ApiSecurityHandlerName).
		WithErrorHandler(errorHandler).
		NotFoundHandler().
		StartServer()
}
