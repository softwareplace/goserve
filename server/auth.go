package server

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/error_handler"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/security/encryptor"
	"log"
	"time"
)

type LoginEntryData struct {
	Username string `json:"username"`
	Password string `json:"encryptor"`
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
	//   - user: An instance of LoginEntryData that contains the username, encryptor, and/or email for user authentication.
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

	// IsValidPassword validates the user-provided plaintext password against the stored encrypted password.
	//
	// This method uses the encryptor package to create a password hash from the provided loginEntryData password
	// and compares it with the encrypted password available in the principal context.
	//
	// A default implementation is available as server.DefaultPasswordValidator[*api_context.DefaultContext]
	//
	// Parameters:
	//   - loginEntryData: The LoginEntryData containing the plaintext password to be validated.
	//   - principal: The principal context of type T, which contains the stored encrypted password.
	//
	// Returns:
	//   - bool: True if the passwords match; false otherwise.
	IsValidPassword(loginEntryData LoginEntryData, principal T) bool
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
			log.Printf("LOGIN/DECRYPT: Failed to decrypt encryptor: %v", err)
		} else {
			loginEntryData.Password = decrypt
		}

		principal, err := loginService.Login(loginEntryData)

		if err != nil {
			log.Printf("LOGIN/LOGIN: Failed to login: %v", err)
			ctx.Forbidden("Login failed: Invalid username or encryptor")
			return
		}

		isValidPassword := loginService.IsValidPassword(loginEntryData, principal)

		if !isValidPassword {
			log.Printf("LOGIN/LOGIN: Failed to login: %v", err)
			ctx.Forbidden("Login failed: Invalid username or encryptor")
			return
		}

		jwt, err := loginService.SecurityService().GenerateJWT(principal, loginService.TokenDuration())

		if err != nil {
			log.Printf("LOGIN/JWT: Failed to generate JWT: %v", err)
			ctx.InternalServerError("Login failed with internal server error. Please try again later.")
			return
		}

		ctx.Ok(jwt)
	}, func(err error) {
		log.Printf("LOGIN/HANDLER: Failed to handle request: %v", err)
		ctx.BadRequest("Login failed: Invalid username or encryptor")
	})
}

// DefaultPasswordValidator is a generic type responsible for validating user passwords
// against their stored encrypted counterparts in principal contexts.
//
// T represents a type that implements the ApiPrincipalContext interface.
// It ensures that the principal context contains methods to retrieve the encrypted password and other details.
type DefaultPasswordValidator[T api_context.ApiPrincipalContext] struct {
}

func (a *DefaultPasswordValidator[T]) IsValidPassword(loginEntryData LoginEntryData, principal T) bool {
	return encryptor.NewEncrypt(loginEntryData.Password).
		IsValidPassword(principal.EncryptedPassword())
}
