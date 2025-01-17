package security

import (
	"github.com/softwareplace/http-utils/api_context"
	"time"
)

type ApiSecurityService[T api_context.ApiContextData] interface {

	// Secret retrieves the secret key used to sign and validate JWT tokens.
	// This function ensures consistent access to the secret key across the service.
	//
	// Returns:
	//   - A byte slice containing the secret key.
	Secret() []byte

	// GenerateApiSecretJWT generates a JWT token with the provided ApiJWTInfo.
	// It encrypts the apiKey provided in jwtInfo, then creates a JWT token with
	// the claims containing the client identifier, encrypted apiKey, and expiration time.
	//
	// Parameters:
	// - jwtInfo: an instance of ApiJWTInfo containing the client identifier, apiKey, and expiration duration.
	//
	// Returns:
	// - string: the signed JWT token on success.
	// - error: an error if token generation fails.
	//
	// Note: This method uses the HS256 signing method and requires a secret key.
	GenerateApiSecretJWT(jwtInfo ApiJWTInfo) (string, error)

	// ExtractJWTClaims validates and extracts the JWT claims from the API request context.
	//
	// This function processes the JWT token provided in the `Authorization` header of the request.
	// It performs the following steps:
	// 1. Parses the JWT token using the configured `Secret` to verify the signature.
	// 2. Extracts claims from the validated token and sets them in the API context if valid.
	// 3. Decrypts the "request" claim to retrieve the requester identity and sets it in the API context.
	//
	// If validation fails at any step (e.g., token parsing, signature verification, claims decoding),
	// it logs the error and sends an appropriate HTTP error response to the client.
	//
	// Parameters:
	// - ctx: The API request context containing the `Authorization` JWT token.
	//
	// Returns:
	// - bool: True if JWT claims are successfully extracted and valid; False otherwise.
	//
	// Notes:
	// - This method relies on the `jwt-go` library for parsing and managing JWT tokens.
	// - Decrypt and cryptographic methods used must ensure secure implementation.
	ExtractJWTClaims(requestContext api_context.ApiRequestContext[T]) bool

	// JWTClaims extracts and parses the claims from the provided JWT token in the ApiRequestContext.
	// It uses the context's ApiKey field as the JWT token for processing. The token is validated
	// and then its claims are extracted into a map of string to interface{}.
	//
	// This function relies on the jwt-go library for processing the token. It uses the Secret()
	// method of the service to provide the secret necessary for verifying the JWT signature.
	//
	// Parameters:
	// - ctx: The ApiRequestContext containing the API key (JWT token to be parsed).
	//
	// Returns:
	// - A map of string to interface{} representing the JWT claims if parsing is successful.
	// - An error if the token parsing fails, the token is invalid, or claims extraction is unsuccessful.
	//
	// Example:
	//   claims, err := securityService.JWTClaims(apiContext)
	//   if err != nil {
	//	   log.Fatalf("failed to extract claims: %v", err)
	//   }
	//   fmt.Printf("Extracted JWT claims: %v", claims)
	JWTClaims(ctx api_context.ApiRequestContext[T]) (map[string]interface{}, error)

	// GenerateJWT
	// apiSecurityServiceImpl provides methods to handle JWT operations such as
	// generation, validation, and extraction of claims. It also deals with secure
	// management of the application's authorization mechanism.
	//
	// Validation is performed using the provided JWT token, ensuring that the
	// claims and tokens are properly extracted, decrypted, and stored within the
	// API context. The service also handles error cases and ensures secure access
	// control across APIs.
	//
	// The jwt-go library is utilized for parsing, validating, and generating JWT tokens.
	//
	// All cryptographic and encoding operations are expected to be implemented
	// within the `Encrypt` and `Decrypt` methods of this service implementation.
	//
	// Usage example for GenerateJWT:
	//
	//	data := YourDataObject{
	//		Salt:  "your_salt",
	//		Roles: []string{"role1", "role2"},
	//	}
	//	securityService := &apiSecurityServiceImpl{} // Properly initialize implementation
	//	tokenDetails, err := securityService.GenerateJWT(data, time.Hour*24)
	//	if err != nil {
	//		log.Fatalf("failed to generate token: %v", err)
	//	}
	//	fmt.Printf("Generated JWT: %v", tokenDetails)
	//
	// The `GenerateJWT` method securely encrypts the `Salt` and `Roles` from the
	// provided data, embeds this information as claims in the token, and generates
	// the JWT that is configured with an expiration time.
	//
	// This service ensures that errors during the JWT process are logged and result
	// in appropriate HTTP error responses.
	//
	// Special Notes:
	// - All cryptographic operations must use secure mechanisms.
	// - Ensure better logging practices during debugging sensitive context data.
	// - It is recommended to configure proper secret rotation policies.
	GenerateJWT(user T, duration time.Duration) (map[string]interface{}, error)

	// Encrypt encrypts the given value using the secret associated with the apiSecurityServiceImpl instance.
	// It returns the encrypted string or an error if encryption fails.
	Encrypt(key string) (string, error)

	// Decrypt decrypts the given encrypted string using the secret associated with the apiSecurityServiceImpl instance.
	// It returns the decrypted string or an error if decryption fails.
	//
	// Parameters:
	// - encrypted: The string that has been encrypted and needs to be decrypted.
	//
	// Returns:
	// - A string representing the decrypted value if the operation is successful.
	// - An error if decryption fails due to issues like invalid cipher text or incorrect secret.
	//
	// Notes:
	// - The decryption logic must use secure cryptographic mechanisms to ensure data safety.
	// - Ensure that any sensitive data involved in the decryption process is handled securely
	//   and not exposed in logs or error messages.
	Decrypt(encrypted string) (string, error)

	// Validation apiSecurityServiceImpl handles API security by providing methods for
	// working with JSON Web Tokens (JWT) such as validation, claim extraction,
	// and token generation.
	//
	// This service ensures secure mechanisms for API authorization and maintains
	// strict validation and processing of JWT tokens using the `jwt-go` library.
	// All cryptographic operations like encryption and decryption are implemented
	// via the Encrypt and Decrypt methods.
	//
	// Key Responsibilities:
	// - Validates JWT tokens for API requests and ensures proper claims extraction.
	// - Generates JWT tokens with custom claims and lifespan.
	// - Handles token signature validation using a configured secret.
	// - Ensures secure access control to the APIs.
	//
	// Notes:
	// - Proper error handling and logging are implemented for fault tolerance.
	// - Secret rotation policies and encryption standards should be followed.
	// - Always ensure sensitive information is never exposed in logs or error messages.
	Validation(
		ctx api_context.ApiRequestContext[T],
		next func(ctx api_context.ApiRequestContext[T]) (*T, bool),
	) (*T, bool)
}

type ApiJWTInfo struct {
	Client string
	Key    string
	// Expiration in hours
	Expiration time.Duration //
}

type apiSecurityServiceImpl[T api_context.ApiContextData] struct {
	ApiSecretAuthorization string
}

var (
	instance apiSecurityServiceImpl[api_context.ApiContextData]
)

func GetApiSecurityService[T api_context.ApiContextData](apiSecretAuthorization string) ApiSecurityService[T] {
	instance.Secret()
	instance := apiSecurityServiceImpl[T]{
		ApiSecretAuthorization: apiSecretAuthorization,
	}
	return &instance
}
