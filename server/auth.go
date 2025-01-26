package server

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/error_handler"
	"github.com/softwareplace/http-utils/security"
	"log"
	"time"
)

type LoginEntryData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type ApiKeyEntryData struct {
	ClientName string        `json:"clientName"` // Client information for which the public key is generated (required)
	Expiration time.Duration `json:"expiration"` // Expiration specifies the duration until the API key expires (optional).
	ClientId   string        `json:"clientId"`   // ClientId represents the unique identifier for the client associated with the API key entry (required).
}

type LoginService[T api_context.ApiPrincipalContext] interface {
	// SecurityService returns an instance of ApiSecurityService responsible for providing security-related operations.
	// This includes handling encryption, decryption, JWT generation, claim extraction, and authorization processes.
	// It ensures the necessary security mechanisms are available in the context of the LoginService.
	//
	// Returns:
	//   - security.ApiSecurityService[T]: The security service instance associated with the implementing service.
	SecurityService() security.ApiSecurityService[T]

	// Login processes the login request for the specified user by validating their credentials.
	// It authenticates the user based on the provided login data and returns an authenticated principal context or an error.
	//
	// Parameters:
	//   - user: An instance of LoginEntryData that contains the username, password, and/or email for user authentication.
	//
	// Returns:
	//   - T: The authenticated principal context representing the logged-in user.
	//   - error: If authentication fails, an error is returned.
	Login(user LoginEntryData) (T, error)

	// TokenDuration specifies the duration for which a generated JWT token remains valid.
	// This value defines the time-to-live (TTL) for the token, ensuring secure and proper session management.
	//
	// Returns:
	//   - time.Duration: The duration for which a generated token is valid.
	TokenDuration() time.Duration
}

type ApiKeyGeneratorService[T api_context.ApiPrincipalContext] interface {

	// SecurityService returns an instance of ApiSecurityService responsible for handling security-related operations.
	// This includes operations such as JWT generation, claims extraction, encryption, decryption, and authorization handling.
	// It provides the foundational security mechanisms required by the ApiKeyGeneratorService.
	//
	// Returns:
	//   - security.ApiSecurityService[T]: The security service instance associated with the implementing service,
	//	 providing security functionalities for API keys, JWTs, and authorization processes.
	SecurityService() security.ApiSecurityService[T]

	// GetApiJWTInfo generates the security.ApiJWTInfo for the given ApiKeyEntryData and ApiRequestContext.
	// This method is responsible for processing the API key entry data and request context to create an ApiJWTInfo object,
	// which contains essential JWT-related information such as the client, key, and expiration details.
	//
	// Parameters:
	//   - apiKeyEntryData: An instance of ApiKeyEntryData that includes client details, expiration duration, and unique client identifier.
	//   - ctx: The API request context, which contains metadata and principal information related to the API key generation process.
	//
	// Returns:
	//   - security.ApiJWTInfo: The generated ApiJWTInfo object containing JWT details necessary for creating the API secret JWT.
	//   - error: If an error occurs during the process, it returns the corresponding error; otherwise, nil.
	GetApiJWTInfo(apiKeyEntryData ApiKeyEntryData, ctx *api_context.ApiRequestContext[T]) (security.ApiJWTInfo, error)

	// OnGenerated is invoked after an API key has been successfully generated.
	// This function allows additional processing or handling, such as logging,
	// auditing, or notifying dependent systems of the newly generated API key.
	//
	// Parameters:
	//   - data: The generated token as security.JwtResponse.
	//   - ctx: The API request context, containing metadata and principal
	//		  information related to the API key generation.
	OnGenerated(data security.JwtResponse, ctx api_context.SampleContext[T])
}

func (a *apiRouterHandlerImpl[T]) Login(ctx *api_context.ApiRequestContext[T]) {
	GetRequestBody(ctx, LoginEntryData{}, a.loginDataHandler, FailedToLoadBody[T])
}

func (a *apiRouterHandlerImpl[T]) ApiKeyGenerator(ctx *api_context.ApiRequestContext[T]) {
	GetRequestBody(ctx, ApiKeyEntryData{}, a.apiKeyGeneratorDataHandler, FailedToLoadBody[T])
}

func (a *apiRouterHandlerImpl[T]) loginDataHandler(ctx *api_context.ApiRequestContext[T], loginEntryData LoginEntryData) {

	error_handler.Handler(func() {
		loginService := a.loginService
		decrypt, err := loginService.SecurityService().Decrypt(loginEntryData.Password)
		if err != nil {
			log.Printf("LOGIN/DECRYPT: Failed to decrypt password: %v", err)
		} else {
			loginEntryData.Password = decrypt
		}

		login, err := loginService.Login(loginEntryData)

		if err != nil {
			log.Printf("LOGIN/LOGIN: Failed to login: %v", err)
			ctx.BadRequest("Login failed: Invalid username or password")
			return
		}

		jwt, err := loginService.SecurityService().GenerateJWT(login, loginService.TokenDuration())

		if err != nil {
			log.Printf("LOGIN/JWT: Failed to generate JWT: %v", err)
			ctx.InternalServerError("Login failed with internal server error. Please try again later.")
			return
		}

		ctx.Ok(jwt)
	}, func(err error) {
		log.Printf("LOGIN/HANDLER: Failed to handle request: %v", err)
		ctx.BadRequest("Login failed: Invalid username or password")
	})
}

func (a *apiRouterHandlerImpl[T]) apiKeyGeneratorDataHandler(ctx *api_context.ApiRequestContext[T], apiKeyEntryData ApiKeyEntryData) {
	error_handler.Handler(func() {
		log.Printf("API/KEY/GENERATOR: requested by: %s", ctx.AccessId)

		jwtInfo, err := a.apiKeyGeneratorService.GetApiJWTInfo(apiKeyEntryData, ctx)

		if err != nil {
			log.Printf("API/KEY/GENERATOR: Failed to generate JWT: %v", err)
			ctx.InternalServerError("Failed to generate JWT. Please try again later.")
			return
		}

		jwt, err := a.loginService.SecurityService().GenerateApiSecretJWT(jwtInfo)

		if err != nil {
			log.Printf("API/KEY/GENERATOR: Failed to generate JWT: %v", err)
			ctx.InternalServerError("Failed to generate JWT. Please try again later.")
			return
		}

		ctx.Ok(jwt)

		a.apiKeyGeneratorService.OnGenerated(*jwt, api_context.SampleCtx[T](*ctx))
	}, func(err error) {
		log.Printf("API/KEY/GENERATOR/HANDLER: Failed to handle request: %v", err)
		ctx.InternalServerError("Failed to generate JWT. Please try again later.")
	})
}
