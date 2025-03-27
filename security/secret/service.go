package secret

import (
	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security"
	"github.com/softwareplace/goserve/security/principal"
)

const (
	AccessHandlerError = "API_SECRET_ACCESS_HANDLER_ERROR"
	AccessHandlerName  = "API_SECRET_MIDDLEWARE"
)

type Service[T goservectx.Principal] interface {
	security.Service[T]
	Provider[T]

	// HandlerSecretAccess is the core function of the Service interface that is responsible for
	// validating API secret keys to ensure secure access to API resources.
	//
	// This method enforces access security by leveraging the following mechanisms:
	//   - Validates the API key against the stored private key to confirm its authenticity.
	//   - Handles any errors that occur during the validation process, responding with appropriate
	//	 HTTP status codes and error messages if validation fails.
	//
	// The function works as follows:
	//   1. Calls the `apiSecretKeyValidation` method to validate the API key and public/private key pair.
	//   2. If the validation fails, invokes the `handlerErrorOrElse` method from the Service to handle
	//	  the failure, typically responding with `http.StatusUnauthorized`.
	//   3. If validation is successful, allows the request to proceed by returning `true`.
	//
	// Args:
	//   - ctx (*context.Request[T]): The context of the incoming API request that carries
	//	 all necessary information for validation, such as JWT claims and keys.
	//
	// Returns:
	//   - bool: `true` if the API key is valid and access is granted; `false` otherwise.
	HandlerSecretAccess(ctx *goservectx.Request[T]) bool

	// DisableForPublicPath sets whether validation should be skipped for public API paths.
	//
	// This method allows configuring the API secret handler to bypass validation for requests
	// targeting public endpoints. When enabled, security mechanisms such as API key validation
	// may not be enforced on these paths, allowing unauthenticated access as needed.
	//
	// Args:
	//   - ignore (bool): A flag indicating whether to ignore validation for public paths.
	//	 Set to `true` to skip validation; set to `false` to enforce validation.
	DisableForPublicPath(ignore bool) Service[T]

	Handler(ctx *goservectx.Request[T], body ApiKeyEntryData)
}

type apiSecretHandlerImpl[T goservectx.Principal] struct {
	security.Service[T]
	Provider[T]
	pService                       principal.Service[T]
	secretKey                      string
	apiSecret                      any
	ignoreValidationForPublicPaths bool
}
