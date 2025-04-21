package security

import (
	log "github.com/sirupsen/logrus"
	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/security/jwt"
	"github.com/softwareplace/goserve/security/principal"
	"github.com/softwareplace/goserve/utils"
)

const (
	ApiSecurityHandlerName = "API_SECURITY_MIDDLEWARE"
)

type JwtResponse struct {
	Token   string `json:"token"`
	Expires int    `json:"expires"`
}

type Service[T goservectx.Principal] interface {
	jwt.Service[T]
	ResourceAccessValidation[T]

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
	AuthorizationHandler(ctx *goservectx.Request[T]) (doNext bool)
}

type impl[T goservectx.Principal] struct {
	ResourceAccessValidation[T]
	jwt.Service[T]
	PService principal.Service[T]
}

// New creates a new instance of the security Service with a default error handler.
//
// This function initializes the Service using the provided API secret authorization key
// and principal service. It also sets up a default resource access handler and error handler.
//
// Parameters:
// - service: The principal service responsible for managing and loading user principals.
//
// Returns:
// - Service[T]: A new instance of the security Service.
func New[T goservectx.Principal](
	service principal.Service[T],
) Service[T] {
	defaultErrorHandler := goserveerror.Default[T]()
	apiSecretKey := utils.GetEnvOrDefault("API_SECRET_KEY", "")

	if apiSecretKey == "" {
		log.Fatal("API_SECRET_KEY environment variable is not set")
	}

	return &impl[T]{
		ResourceAccessValidation: &defaultResourceAccessHandler[T]{
			&defaultErrorHandler,
		},
		Service:  jwt.New(service, apiSecretKey, defaultErrorHandler),
		PService: service,
	}
}

// Create creates a new instance of the security Service with the provided configurations.
//
// This function is a more customizable version of New where you can provide your own error
// handler and resource access validation logic.
//
// Parameters:
//   - apiSecretKey: The secret key used for API authorization and JWT management, encrypt and decrypt values.
//   - service: The principal service responsible for managing and loading user principals.
//   - handler: A pointer to a custom API error handler that processes authorization errors.
//   - resourceValidation: A custom resource access validation implementation.
//
// Returns:
// - Service[T]: A new instance of the security Service with the provided configurations.
func Create[T goservectx.Principal](
	apiSecretKey string,
	service principal.Service[T],
	handler goservectx.ApiHandler[T],
	resourceValidation ResourceAccessValidation[T],
) Service[T] {
	return &impl[T]{
		ResourceAccessValidation: resourceValidation,
		Service:                  jwt.New(service, apiSecretKey, handler),
		PService:                 service,
	}
}
