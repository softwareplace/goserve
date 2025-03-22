package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/error_handler"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/security/encryptor"
	"github.com/softwareplace/http-utils/security/principal"
	"github.com/softwareplace/http-utils/server"
	"github.com/softwareplace/http-utils/test/gen"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type loginServiceImpl struct {
	server.DefaultPasswordValidator[*api_context.DefaultContext]
	securityService security.ApiSecurityService[*api_context.DefaultContext]
}

var mockStore = map[string]string{
	"37c75552614a4eb58a2eb2d04928cdfd": "D/b7o5KGWe0SOF06r7bvKWyud95XVQwD9xp9NIDqMUWqt1xHz6PpIAF2jRo6pFGaaTwglXwql7QChU1fmQf7omQnjZImS9iWhKh9xvQEpXhygA5WAzBEPiekmyfH6LwkWgFQeFxi4spwX5J+m1LPMIrHZyjVqFOr01f3RaHAlBwxOwWdbQ0au32gVshGFY7Rt7d5RmMQATA0rQf0NGZlcIEM5ez8hBxjUHnKakGjYOITQsd570wvlFnRhvkvoxRfpAGAexXRAS8tImdiw/L7BVSbTKjwqSfweH59CK3JhHC/qdwDlSDA6rJWat4MOeb2qWbgbmlQV71QEFOZ9k78gdNz3FuFsMIQ4Swyf3dvBraTFlCjxDil7fIyTT1PJ8f8AvMcVdzWsXwWRl5+SgJvHcZI9nGmswzacRv2T008qUKm28m6By5Sd1ux38QghobBtpL2n3+lgEnov59/cStPHS4kSNrudeX1RtU7DPlqWZUyXkn4H+3tdlUXMufZcYekIkq3fIVsGHxRRGTRA1ILell9FBXwEVw/je2FsrzIZbPxZKnRb8WRbqNFreDf/9hdWLjKw4IaIddRUbGUSTLV3u94QbhDwsdFRmorMgKZd3yukVc=",
	/// X-Api-Key for test only eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcGlLZXkiOiJTL2VTYzVDQ3Jub1laaDAyU2pLdFVsSzFXdmRaaVA1OXpFUU9jNE54K0pjL1c1dkhMa0tndE1ueExHN3dKTUwvIiwiY2xpZW50IjoiU29mdHdhcmUgUGxhY2UiLCJleHAiOjMwMzM5MzczNTcsInNjb3BlIjpbIkhFMSs0cEVwM3YzZFBzWXNLa3FLMGkzdiswSjMvYjFVN01YQkx3ZzhxQ0E9IiwiR2lQWUVNU1IvK1BjNUdaTm9OcUpqZDRkS1FZbjZ6QzBMbmdYTHVxdFc4VzkiLCJjY294TWNaT0tEZ0srTUZuend0YWFEWXgxaEtPSVlKNDl3PT0iLCJxOWRHb3V5bTBxZWxvV1V4bElKZ2Y1U3l6UnIrU3YwWWwvVT0iLCJNOHdrRkN3cmZpeVBKc2hjb3NrQU5GS0RZZ2ZxRnJOWXkwVmljOEdlM3dPSyJdfQ.n5_8kp3nNqXOAZVB73GCIXcv61gNyyihqz6xDIjIA0k
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

func (l *loginServiceImpl) RequiredScopes() []string {
	return []string{
		"api:key:generator",
	}
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
	apiJWTInfo security.ApiJWTInfo,
	ctx api_context.SampleContext[*api_context.DefaultContext],
) {
	mockStore[apiJWTInfo.Key] = *apiJWTInfo.PublicKey
	log.Printf("%s - %s", apiJWTInfo.Key, data.Token)
	log.Printf("API KEY GENERATED: from %s - %v", ctx.AccessId, data)
}

func (l *loginServiceImpl) Login(user server.LoginEntryData) (*api_context.DefaultContext, error) {
	result := &api_context.DefaultContext{}
	result.SetRoles("api:example:user", "api:example:admin", "read:pets", "write:pets", "api:key:generator")
	password := encryptor.NewEncrypt(user.Password).EncodedPassword()
	result.SetEncryptedPassword(password)
	return result, nil
}

func (l *loginServiceImpl) TokenDuration() time.Duration {
	return time.Minute * 15
}

type secretProviderImpl []struct{}

func (s *secretProviderImpl) Get(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) (string, error) {
	return mockStore[ctx.ApiKeyId], nil
}

type principalServiceImpl struct {
}

func (d *principalServiceImpl) LoadPrincipal(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) bool {
	if ctx.Authorization == "" {
		return false

	}

	context := api_context.NewDefaultCtx()
	context.SetRoles("api:key:generator")
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

var (
	userPrincipalService principal.PService[*api_context.DefaultContext]            = &principalServiceImpl{}
	errorHandler         error_handler.ApiErrorHandler[*api_context.DefaultContext] = &errorHandlerImpl{}
	securityService                                                                 = security.ApiSecurityServiceBuild(
		"ue1pUOtCGaYS7Z1DLJ80nFtZ",
		userPrincipalService,
	)

	loginService = &loginServiceImpl{
		securityService: securityService,
	}

	secretProvider = secretProviderImpl{}
	secretHandler  = security.ApiSecretAccessHandlerBuild(
		"./secret/private.key",
		&secretProvider,
		securityService,
	)

	apiSecret = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcGlLZXkiOiJTL2VTYzVDQ3Jub1laaDAyU2pLdFVsSzFXdmRaaVA1OXpFUU9jNE54K0pjL1c1dkhMa0tndE1ueExHN3dKTUwvIiwiY2xpZW50IjoiU29mdHdhcmUgUGxhY2UiLCJleHAiOjMwMzM5MzczNTcsInNjb3BlIjpbIkhFMSs0cEVwM3YzZFBzWXNLa3FLMGkzdiswSjMvYjFVN01YQkx3ZzhxQ0E9IiwiR2lQWUVNU1IvK1BjNUdaTm9OcUpqZDRkS1FZbjZ6QzBMbmdYTHVxdFc4VzkiLCJjY294TWNaT0tEZ0srTUZuend0YWFEWXgxaEtPSVlKNDl3PT0iLCJxOWRHb3V5bTBxZWxvV1V4bElKZ2Y1U3l6UnIrU3YwWWwvVT0iLCJNOHdrRkN3cmZpeVBKc2hjb3NrQU5GS0RZZ2ZxRnJOWXkwVmljOEdlM3dPSyJdfQ.n5_8kp3nNqXOAZVB73GCIXcv61gNyyihqz6xDIjIA0k"
)

func TestMockServer(t *testing.T) {
	t.Run("expects that can get login response successfully", func(t *testing.T) {
		// Create a new request
		loginBody := strings.NewReader(`{"username": "my-username","password": "ynT9558iiMga&ayTVGs3Gc6ug1"}`)
		req, err := http.NewRequest("POST", "/login", loginBody)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		server.Default().
			WithLoginResource(loginService).
			WithApiSecurityService(securityService).
			WithPrincipalService(userPrincipalService).
			NotFoundHandler().
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("expects that return 401 when api secret is required for all resources but was not provided", func(t *testing.T) {
		// Create a new request
		loginBody := strings.NewReader(`{"username": "my-username","password": "ynT9558iiMga&ayTVGs3Gc6ug1"}`)
		req, err := http.NewRequest("POST", "/login", loginBody)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		secretProvider := secretProviderImpl{}
		secretHandler := security.ApiSecretAccessHandlerBuild(
			"./secret/private.key",
			&secretProvider,
			securityService,
		)

		server.Default().
			WithLoginResource(loginService).
			WithApiSecretAccessHandler(secretHandler).
			WithApiSecurityService(securityService).
			WithPrincipalService(userPrincipalService).
			NotFoundHandler().
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusUnauthorized)
		}
	})

	t.Run("expects that can get login response successfully when requires api secret and it was provided", func(t *testing.T) {
		// Create a new request
		loginBody := strings.NewReader(`{"username": "my-username","password": "ynT9558iiMga&ayTVGs3Gc6ug1"}`)
		req, err := http.NewRequest("POST", "/login", loginBody)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set(api_context.XApiKey, apiSecret)

		rr := httptest.NewRecorder()

		server.Default().
			WithLoginResource(loginService).
			WithApiSecretAccessHandler(secretHandler).
			WithApiSecurityService(securityService).
			WithPrincipalService(userPrincipalService).
			NotFoundHandler().
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("expects that return default not found when a custom was not provided", func(t *testing.T) {
		// Create a new request
		req, err := http.NewRequest("POST", "/not-found", nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		server.Default().
			WithPrincipalService(userPrincipalService).
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}

		if strings.Contains(rr.Body.String(), "404 page not found") {
			t.Log("Response body contains '404 page not found'")
		} else {
			t.Errorf("Expected response body to contain '404 page not found', but got: %s", rr.Body.String())
		}

	})

	t.Run("expects that return custom not found when a custom was provided", func(t *testing.T) {
		// Create a new request
		req, err := http.NewRequest("POST", "/not-found", nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		server.Default().
			WithPrincipalService(userPrincipalService).
			CustomNotFoundHandler(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte("Custom 404 Page"))
			}).
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}

		if strings.Contains(rr.Body.String(), "Custom 404 Page") {
			t.Log("Response body contains 'Custom 404 Page'")
		} else {
			t.Errorf("Expected response body to contain 'Custom 404 Page', but got: %s", rr.Body.String())
		}
	})

	t.Run("expects that return swagger resource when swagger was defined and using the default not found handler", func(t *testing.T) {
		// Create a new request
		req, err := http.NewRequest("GET", "/", nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		server.Default().
			WithPrincipalService(userPrincipalService).
			EmbeddedServer(gen.ApiResourceHandler(&_service{})).
			SwaggerDocHandler("./resource/pet-store.yaml").
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMovedPermanently {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusMovedPermanently)
		}

		if strings.Contains(rr.Body.String(), "<a href=\"/swagger/index.html\">Moved Permanently</a>.") {
			t.Log("Response body contains '<a href=\"/swagger/index.html\">Moved Permanently</a>.'")
		} else {
			t.Errorf("Expected response body to contain '<a href=\"/swagger/index.html\">Moved Permanently</a>.', but got: %s", rr.Body.String())
		}
	})
}
