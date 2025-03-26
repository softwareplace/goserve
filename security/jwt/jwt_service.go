package jwt

import (
	apicontext "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security/encryptor"
	"github.com/softwareplace/goserve/security/principal"
	"time"
)

type Entry struct {
	Client     string
	Key        string
	Expiration time.Duration
	Scopes     []string
	PublicKey  *string
}

type Response struct {
	Token   string `json:"token"`
	Expires int    `json:"expires"`
}

type Service[T apicontext.Principal] interface {
	encryptor.Service
	// GenerateApiSecretJWT generates a JWT token with the provided Entry.
	// It encrypts the apiKey provided in jwtInfo, then creates a JWT token with
	// the claims containing the client identifier, encrypted apiKey, and expiration time.
	//
	// Parameters:
	// - entry: an instance of Entry containing the client identifier, apiKey, and expiration duration.
	//
	// Returns:
	// - string: the signed JWT token on success.
	// - error: an error if token generation fails.
	//
	// Note: This method uses the HS256 signing method and requires a secret key.
	GenerateApiSecretJWT(entry Entry) (*Response, error)

	// ExtractJWTClaims validates and extracts the JWT claims from the API request context.
	//
	// This function processes the JWT token provided in the `AuthorizationHandler` header of the request.
	// It performs the following steps:
	// 1. Parses the JWT token using the configured `Secret` to verify the signature.
	// 2. Extracts claims from the validated token and sets them in the API context if valid.
	// 3. Decrypts the "request" claim to retrieve the requester identity and sets it in the API context.
	//
	// If validation fails at any step (e.g., token parsing, signature verification, claims decoding),
	// it logs the error and sends an appropriate HTTP error response to the client.
	//
	// Parameters:
	// - ctx: The API request context containing the `AuthorizationHandler` JWT token.
	//
	// Returns:
	// - bool: True if JWT claims are successfully extracted and valid; False otherwise.
	//
	// Notes:
	// - This method relies on the `jwt-go` library for parsing and managing JWT tokens.
	// - Decrypt and cryptographic methods used must ensure secure implementation.
	ExtractJWTClaims(requestContext *apicontext.Request[T]) bool

	// JWTClaims extracts and parses the claims from the provided JWT token in the Request.
	// It uses the context's ApiKey field as the JWT token for processing. The token is validated
	// and then its claims are extracted into a map of string to interface{}.
	//
	// This function relies on the jwt-go library for processing the token. It uses the Secret()
	// method of the PService to provide the secret necessary for verifying the JWT signature.
	//
	// Parameters:
	// - ctx: The Request containing the API key (JWT token to be parsed).
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
	JWTClaims(ctx *apicontext.Request[T]) (map[string]interface{}, error)

	// GenerateJWT
	// apiSecurityServiceImpl provides methods to handle JWT operations such as
	// generation, validation, and extraction of claims. It also deals with secure
	// management of the application's authorization mechanism.
	//
	// Validation is performed using the provided JWT token, ensuring that the
	// claims and tokens are properly extracted, decrypted, and stored within the
	// API context. The PService also handles error cases and ensures secure access
	// control across APIs.
	//
	// The jwt-go library is utilized for parsing, validating, and generating JWT tokens.
	//
	// All cryptographic and encoding operations are expected to be implemented
	// within the `Encrypt` and `Decrypt` methods of this PService implementation.
	//
	// Usage example for GenerateJWT:
	//
	//	data := YourDataObject{
	//		GetId:  "your_salt",
	//		GetRoles: []string{"role1", "role2"},
	//	}
	//	securityService := &apiSecurityServiceImpl{} // Properly initialize implementation
	//	tokenDetails, err := securityService.GenerateJWT(data, time.Hour*24)
	//	if err != nil {
	//		log.Fatalf("failed to generate token: %v", err)
	//	}
	//	fmt.Printf("Generated JWT: %v", tokenDetails)
	//
	// The `GenerateJWT` method securely encrypts the `GetId` and `GetRoles` from the
	// provided data, embeds this information as claims in the token, and generates
	// the JWT that is configured with an expiration time.
	//
	// This PService ensures that errors during the JWT process are logged and result
	// in appropriate HTTP error responses.
	//
	// Special Notes:
	// - All cryptographic operations must use secure mechanisms.
	// - Ensure better logging practices during debugging sensitive context data.
	// - It is recommended to configure proper secret rotation policies.
	GenerateJWT(user T, duration time.Duration) (*Response, error)

	HandlerErrorOrElse(
		ctx *apicontext.Request[T],
		error error,
		executionContext string,
		handlerNotFound func(),
	)
}

type serviceImpl[T apicontext.Principal] struct {
	encryptor.Service
	PService     principal.Service[T]
	ErrorHandler apicontext.ApiHandler[T]
}

func New[T apicontext.Principal](
	pService principal.Service[T],
	apiSecretAuthorization string,
	errorHandler apicontext.ApiHandler[T],
) Service[T] {
	return &serviceImpl[T]{
		Service:      encryptor.New([]byte(apiSecretAuthorization)),
		PService:     pService,
		ErrorHandler: errorHandler,
	}
}
