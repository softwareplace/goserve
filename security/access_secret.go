package security

import (
	"github.com/softwareplace/http-utils/apicontext"
)

const (
	ApiSecretAccessHandlerError = "API_SECRET_ACCESS_HANDLER_ERROR"
	ApiSecretAccessHandlerName  = "API_SECRET_MIDDLEWARE"
)

// ApiSecretKeyProvider is an interface designed to provide secure access to API secret keys
// based on the context of an incoming API request.
//
// This interface plays a crucial role in enabling secure communication and access control
// by retrieving API keys that are specific to each request. It abstracts the mechanism for
// obtaining these keys, which could involve fetching from a database, a configuration file,
// or any other secure storage mechanism.
//
// Type Parameters:
//   - T: A type that satisfies the `apicontext.ApiPrincipalContext` interface, representing
//     the authentication and authorization context for API requests.
type ApiSecretKeyProvider[T apicontext.ApiPrincipalContext] interface {

	// Get (ctx *apicontext.ApiRequestContext[T]) (string, error):
	//	   Fetches the API secret key for the given request context. The method should implement
	//	   any necessary logic to securely retrieve and provide the key, such as decryption or
	//	   validation.
	//
	// Example Use Case:
	// When processing an API request that requires validation with a secret key, the implementation
	// of this interface can retrieve and provide the appropriate key tailored to the request's context.
	//
	// Returns:
	//   - A string representing the API secret key.
	//   - An error if the key retrieval or processing fails, ensuring proper error handling in the
	//	 request lifecycle.
	Get(ctx *apicontext.ApiRequestContext[T]) (string, error)
}

type ApiSecretAccessHandler[T apicontext.ApiPrincipalContext] interface {

	// HandlerSecretAccess is the core function of the ApiSecretAccessHandler interface that is responsible for
	// validating API secret keys to ensure secure access to API resources.
	//
	// This method enforces access security by leveraging the following mechanisms:
	//   - Validates the API key against the stored private key to confirm its authenticity.
	//   - Handles any errors that occur during the validation process, responding with appropriate
	//	 HTTP status codes and error messages if validation fails.
	//
	// The function works as follows:
	//   1. Calls the `apiSecretKeyValidation` method to validate the API key and public/private key pair.
	//   2. If the validation fails, invokes the `handlerErrorOrElse` method from the ApiSecurityService to handle
	//	  the failure, typically responding with `http.StatusUnauthorized`.
	//   3. If validation is successful, allows the request to proceed by returning `true`.
	//
	// Args:
	//   - ctx (*apicontext.ApiRequestContext[T]): The context of the incoming API request that carries
	//	 all necessary information for validation, such as JWT claims and keys.
	//
	// Returns:
	//   - bool: `true` if the API key is valid and access is granted; `false` otherwise.
	HandlerSecretAccess(ctx *apicontext.ApiRequestContext[T]) bool

	// DisableForPublicPath sets whether validation should be skipped for public API paths.
	//
	// This method allows configuring the API secret handler to bypass validation for requests
	// targeting public endpoints. When enabled, security mechanisms such as API key validation
	// may not be enforced on these paths, allowing unauthenticated access as needed.
	//
	// Args:
	//   - ignore (bool): A flag indicating whether to ignore validation for public paths.
	//	 Set to `true` to skip validation; set to `false` to enforce validation.
	DisableForPublicPath(ignore bool) ApiSecretAccessHandler[T]

	//
	// SecretKey provides a secure mechanism to fetch and return the current secret key used for
	// API validations or other security-related operations.
	//
	// This method is essential in ensuring that security-critical processes access the correct
	// API secret key stored or configured in the implementation of the `ApiSecretAccessHandler` interface.
	//
	// Returns:
	//   - string: The current secret key being utilized for authentication or validation purposes.
	SecretKey() string
}
