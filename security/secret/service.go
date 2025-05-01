package secret

import (
	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security"
	"github.com/softwareplace/goserve/security/model"
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
	//   1. Calls the `ApiSecretKeyValidation` method to validate the API key and public/private key pair.
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

	// Handler processes API key generation requests and returns a JWT response.
	// It validates the API key entry data, generates a public key if not present, creates
	// the JWT, and sends the response back to the client.
	//
	// Parameters:
	//   - ctx: The request context which includes request-related metadata and the ability to send responses.
	//   - apiKeyEntryData: The API key entry data which contains information such as client name, expiration duration,
	//	 and a unique client ID.
	//
	// Behavior:
	//   - Logs the API key generation request.
	//   - Retrieves JWT entry information based on the provided API key data.
	//   - If the public key is missing, generates a public key for the request.
	//   - Constructs the JWT token with roles and expiration information.
	//   - Sends the generated JWT as a response to the client.
	//   - Logs and sends an internal server error response if any operation fails.
	//
	// Errors:
	//   - If any step fails (e.g., JWT generation, public key retrieval), the function logs the error and
	//	 responds with an internal server error to the client.
	Handler(ctx *goservectx.Request[T], body model.ApiKeyEntryData)

	// SecretKey returns the current API secret key used for cryptographic operations and token generation.
	SecretKey() string

	// InitAPISecretKey initializes the API secret key by reading and parsing a private key file.
	//
	// This function performs the following steps:
	// - Reads the private key file specified by the `secretKey` field.
	// - Decodes the PEM block from the file data.
	// - Parses the private key from the PEM data using PKCS8 format.
	// - Validates the type of the private key (either ECDSA or RSA).
	// - Stores the private key in the `apiSecret` field of the `apiSecretHandlerImpl` struct.
	//
	// Logs an error and terminates the application if any of the above steps fail.
	InitAPISecretKey()

	// ApiSecretKeyValidation verifies the validity of a public key against the private key stored in the handler.
	//
	// This function performs the following steps:
	// - Extracts JWT claims from the request context using the Service.
	// - Loads the API secret using the provided `ApiSecretKeyServiceProvider`.
	// - Decrypts the API access key to retrieve the PEM-encoded public key.
	// - Decodes the PEM-encoded public key and parses it into a usable public key object.
	// - Validates the type of the parsed public key (ECDSA or RSA).
	// - Ensures the extracted public key corresponds to the private key stored in the `apiSecret` field.
	//
	// If any of the above steps fail, the function logs the error and returns `false`, indicating that
	// the public key validation has failed. Otherwise, it returns `true`.
	//
	// Args:
	//
	//	  ctx (*api_context.Request[T]):
	//		 - The context of the API request carrying the necessary data for validation.
	//
	// Returns:
	//
	//	  bool:
	//		 - `true` if the public key is valid and corresponds to the private key.
	//		 - `false` if the public key is invalid or the validation fails.
	ApiSecretKeyValidation(ctx *goservectx.Request[T]) bool

	// GeneratePubKey generates an encrypted public key from a given private key file.
	//
	// This function performs the following steps:
	// - Reads the private key from the specified file path.
	// - Decodes the PEM block from the private key data.
	// - Parses the private key using the PKCS8 format.
	// - Determines the type of the private key (ECDSA or RSA).
	// - Marshals the corresponding public key into PEM format.
	// - Encrypts the generated PEM-encoded public key using the securityService's encryption logic.
	//
	// Arguments:
	//   - secretKey (string): The file path to the private key.
	//
	// Returns:
	//   - (string, error): An encrypted PEM-encoded public key and an error (if any occurred).
	//
	// Errors:
	//   - Fails if the private key file cannot be read, parsed, or if the key type is unsupported.
	//   - Fails if the public key cannot be marshaled or encrypted.
	//
	// Example:
	//
	//	 encryptedPubKey, err := handler.generatePubKey("path/to/private.key")
	//	 if err != nil {
	//		 log.Printf("Error generating public key: %v", err)
	//	 }
	GeneratePubKey(secretKey string) (string, error)
}

type apiSecretHandlerImpl[T goservectx.Principal] struct {
	security.Service[T]
	Provider[T]
	pService                       principal.Service[T]
	secretKey                      string
	apiSecret                      any
	ignoreValidationForPublicPaths bool
}
