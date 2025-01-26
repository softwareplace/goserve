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
