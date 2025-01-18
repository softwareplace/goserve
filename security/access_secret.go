package security

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/security/principal"
	"log"
	"net/http"
	"os"
)

const (
	ApiSecretAccessHandlerError = "API_SECRET_ACCESS_HANDLER_ERROR"
	ApiSecretAccessHandlerName  = "API_SECRET_MIDDLEWARE"
)

type ApiSecretKeyLoader[T api_context.ApiPrincipalContext] func(ctx *api_context.ApiRequestContext[T]) (string, error)

type ApiSecretAccessHandler[T api_context.ApiPrincipalContext] interface {

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
	//   - ctx (*api_context.ApiRequestContext[T]): The context of the incoming API request that carries
	//	 all necessary information for validation, such as JWT claims and keys.
	//
	// Returns:
	//   - bool: `true` if the API key is valid and access is granted; `false` otherwise.
	HandlerSecretAccess(ctx *api_context.ApiRequestContext[T]) bool

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
}

type apiSecretHandlerImpl[T api_context.ApiPrincipalContext] struct {
	secretKey                      string
	service                        ApiSecurityService[T]
	apiSecret                      any
	loader                         ApiSecretKeyLoader[T]
	principalService               *principal.PService[T]
	ignoreValidationForPublicPaths bool
}

// ApiSecretAccessHandlerBuild apiSecretHandlerImpl is an implementation of the ApiSecretAccessHandler interface which manages
// security-related operations for API requests, such as validating API keys and initializing
// cryptographic keys. It encapsulates the logic for validating an API secret key and restricting
// unauthorized access to resources.
//
// Type Parameters:
//   - T: A type that satisfies the `api_context.ApiPrincipalContext` interface, providing API principal-specific context.
//
// Fields:
//   - secretKey: The file path to the secret key used for cryptographic operations.
//   - service: An instance of ApiSecurityService responsible for cryptographic and security services.
//   - apiSecret: Holder of the parsed private key, supporting either ECDSA or RSA key types.
//   - loader: A function responsible for loading the API secret key for access validation.
//   - principalService: A service managing API principal claims and IDs to ensure request security.
//
// This struct provides methods to initialize the secret key, validate the public key against the private key,
// and enforce access security middleware, ensuring requests are authorized with proper credentials.
func ApiSecretAccessHandlerBuild[T api_context.ApiPrincipalContext](
	secretKey string,
	loader ApiSecretKeyLoader[T],
	service ApiSecurityService[T],
) ApiSecretAccessHandler[T] {
	handler := &apiSecretHandlerImpl[T]{
		secretKey: secretKey,
		service:   service,
		loader:    loader,
	}
	handler.initAPISecretKey()
	return handler
}

func (a *apiSecretHandlerImpl[T]) DisableForPublicPath(ignore bool) ApiSecretAccessHandler[T] {
	a.ignoreValidationForPublicPaths = ignore
	return a
}

func (a *apiSecretHandlerImpl[T]) HandlerSecretAccess(ctx *api_context.ApiRequestContext[T]) bool {
	if a.ignoreValidationForPublicPaths && principal.IsPublicPath[T](*ctx) {
		return true
	}

	if !a.apiSecretKeyValidation(ctx) {
		a.service.handlerErrorOrElse(ctx, nil, ApiSecretAccessHandlerError, func() {
			ctx.Error("You are not allowed to access this resource", http.StatusUnauthorized)
		})
		return false
	}
	return true
}

// initAPISecretKey initializes the API secret key by reading and parsing a private key file.
//
// This function performs the following steps:
// - Reads the private key file specified by the `secretKey` field.
// - Decodes the PEM block from the file data.
// - Parses the private key from the PEM data using PKCS8 format.
// - Validates the type of the private key (either ECDSA or RSA).
// - Stores the private key in the `apiSecret` field of the `apiSecretHandlerImpl` struct.
//
// Logs an error and terminates the application if any of the above steps fail.
func (a *apiSecretHandlerImpl[T]) initAPISecretKey() {
	// Load private key from the provided secretKey file path
	privateKeyData, err := os.ReadFile(a.secretKey)
	if err != nil {
		log.Fatalf("Failed to read private key file: %s", err.Error())
	}

	// Decode PEM block from the private key data
	block, _ := pem.Decode(privateKeyData)
	if block == nil || block.Type != "PRIVATE KEY" {
		log.Fatalf("Failed to decode private key PEM block")
	}

	// Parse the private key using ParsePKCS8PrivateKey
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %s", err.Error())
	}
	a.apiSecret = privateKey

	switch key := a.apiSecret.(type) {
	case *ecdsa.PrivateKey:
		log.Println("Loaded ECDSA private key successfully")
	case *rsa.PrivateKey:
		log.Println("Loaded RSA private key successfully")
	default:
		log.Fatalf("Unsupported private key type: %T", key)
	}
}

// apiSecretKeyValidation verifies the validity of a public key against the private key stored in the handler.
//
// This function performs the following steps:
// - Extracts JWT claims from the request context using the ApiSecurityService.
// - Loads the API secret using the provided `ApiSecretKeyLoader`.
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
//	  ctx (*api_context.ApiRequestContext[T]):
//		 - The context of the API request carrying the necessary data for validation.
//
// Returns:
//
//	  bool:
//		 - `true` if the public key is valid and corresponds to the private key.
//		 - `false` if the public key is invalid or the validation fails.
func (a *apiSecretHandlerImpl[T]) apiSecretKeyValidation(ctx *api_context.ApiRequestContext[T]) bool {
	// Decode the Base64-encoded public key
	claims, err := a.service.JWTClaims(*ctx)

	if err != nil {
		log.Printf("JWT/CLAIMS_EXTRACT: AuthorizationHandler failed: %v", err)
		return false
	}

	(*a.principalService).SetApiKeyClaims(claims)

	apiKey, err := a.service.Decrypt(claims["apiKey"].(string))

	if err != nil {
		log.Printf("JWT/CLAIMS_EXTRACT: AuthorizationHandler failed: %v", err)
		return false
	}

	(*a.principalService).SetApiKeyId(apiKey)

	apiAccessKey, err := a.loader(ctx)
	if err != nil {
		log.Printf("API_SECRET_LOADER: AuthorizationHandler failed: %v", err)
		return false
	}

	// Decode the PEM-encoded public key
	decryptKey, err := a.service.Decrypt(apiAccessKey)
	if err != nil {
		log.Printf("API_SECRET_DECRYPT: AuthorizationHandler failed: %v", err)
		return false
	}
	block, _ := pem.Decode([]byte(decryptKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Printf("Failed to decode public key PEM block")
		return false
	}

	// Parse the public key
	parsedPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Printf("Failed to parse public key: %v", err)
		return false
	}

	switch privateKey := a.apiSecret.(type) {
	case *ecdsa.PrivateKey:
		// Ensure the type of the public key matches ECDSA
		publicKey, ok := parsedPublicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Printf("invalid public key type, expected ECDSA")
			return false
		}

		// Validate if the public key corresponds to the private key
		privateKeyPubKey := &privateKey.PublicKey
		if publicKey.X.Cmp(privateKeyPubKey.X) != 0 || publicKey.Y.Cmp(privateKeyPubKey.Y) != 0 {
			log.Printf("public key does not match the private key")
			return false
		}
	case *rsa.PrivateKey:
		// Ensure the type of the public key matches RSA
		publicKey, ok := parsedPublicKey.(*rsa.PublicKey)
		if !ok {
			log.Printf("invalid public key type, expected RSA")
			return false
		}

		// Validate if the public key corresponds to the private key
		privateKeyPubKey := &privateKey.PublicKey
		if publicKey.E != privateKeyPubKey.E || publicKey.N.Cmp(privateKeyPubKey.N) != 0 {
			log.Printf("public key does not match the private key")
			return false
		}
	default:
		log.Printf("unsupported private key type: %T", privateKey)
		return false
	}

	return true
}
