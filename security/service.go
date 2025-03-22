package security

import (
	apicontext "github.com/softwareplace/http-utils/context"
	errorhandler "github.com/softwareplace/http-utils/error"
	"github.com/softwareplace/http-utils/security/jwt"
	"github.com/softwareplace/http-utils/security/principal"
)

const (
	ApiSecurityHandlerName = "API_SECURITY_MIDDLEWARE"
)

type JwtResponse struct {
	Token   string `json:"token"`
	Expires int    `json:"expires"`
}

type Service[T apicontext.Principal] interface {
	jwt.Service[T]

	// AuthorizationHandler
	// This method is invoked to handle API requests and manage security validation processes.
	// It determines whether the request can proceed further (doNext) based on:
	// 1. Whether the request is made to a public path.
	// 2. The success of the JWT token validation process, which involves:
	//   - Principal extraction.
	//   - Validation of token claims.
	//   - Ensuring proper API authorization.
	//
	// Parameters:
	// - ctx: The Request containing the context information for the API request.
	//
	// Returns:
	// - `true` (doNext) if the request is allowed to continue processing.
	// - `false` if the request fails validation or is unauthorized.
	//
	// Notes:
	// - This function leverages methods like Validation and IsPublicPath to make security decisions.
	// - Ensure that all sensitive operations and data are securely processed.
	// - Public paths bypass validation by default, so it's critical to properly define such paths to avoid security issues.
	AuthorizationHandler(ctx *apicontext.Request[T]) (doNext bool)
}

type impl[T apicontext.Principal] struct {
	jwt.Service[T]
	PService principal.Service[T]
}

func New[T apicontext.Principal](
	apiSecretAuthorization string,
	service principal.Service[T],
	errorHandler errorhandler.ApiErrorHandler[T],
) Service[T] {
	return &impl[T]{
		jwt.New(service, apiSecretAuthorization, errorHandler),
		service,
	}
}

func (a *impl[T]) AuthorizationHandler(ctx *apicontext.Request[T]) (doNext bool) {
	a.ExtractJWTClaims(ctx)
	if principal.IsPublicPath[T](*ctx) {
		return true
	}
	return a.PService.LoadPrincipal(ctx)
}
