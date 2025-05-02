package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security/encryptor"
	"github.com/softwareplace/goserve/security/jwt/model"
	"github.com/softwareplace/goserve/security/principal"
	"time"
)

const (
	IAT = "iat"
	EXP = "exp"
	AUD = "aud"
	SUB = "sub"
	ISS = "iss"
)

type Service[T goservectx.Principal] interface {
	encryptor.Service

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
	ExtractJWTClaims(requestContext *goservectx.Request[T]) bool

	// GetClaims extracts and validates the JWT claims from the provided token.
	//
	// Parameters:
	//   - token: The JWT token from which claims are to be extracted.
	//
	// Returns:
	//   - jwt.MapClaims: A map containing the claims extracted from the token.
	//   - bool: A boolean value indicating whether the claims were successfully extracted and valid.
	GetClaims(token *jwt.Token) (jwt.MapClaims, bool)

	// Generate creates a new JWT Response based on the provided user and duration.
	// It returns the generated Response or an error if the process fails.
	Generate(user T, duration time.Duration) (*model.Response, error)

	// From generates a Response containing a JWT for the given subject, roles, and duration or returns an error if it fails.
	From(sub string, roles []string, duration time.Duration) (*model.Response, error)

	// Issuer returns the identifier of the entity responsible for issuing the JWT tokens in the service.
	Issuer() string

	// Decode extracts claims from a given JWT token string.
	// It validates the token and parses its claims into a map[string]interface{}.
	// Returns an error if the token is invalid or if the claims' structure is incorrect.
	Decode(tokenString string) (map[string]interface{}, error)

	// Parse parses and validates a JWT token string.
	// It uses the secret key provided by the BaseService for token signing and validation.
	// Returns the parsed *jwt.Token or an error if the token cannot be parsed or is invalid.
	Parse(tokenString string) (*jwt.Token, error)

	// IsValid checks if the provided JWT token string is valid.
	// It parses the token string using the configured secret key and verifies the token's validity.
	//
	// Parameters:
	//   - tokenString: The JWT token string to be validated.
	//
	// Returns:
	//   - True if the token is successfully parsed and is valid; otherwise, false.
	IsValid(tokenString string) bool

	HandlerErrorOrElse(
		ctx *goservectx.Request[T],
		error error,
		executionContext string,
		handlerNotFound func(),
	)
}

type BaseService[T goservectx.Principal] struct {
	encryptor.Service
	PService        principal.Service[T]
	ErrorHandler    goservectx.ApiHandler[T]
	claimsExtractor claimsExtractor
}

func New[T goservectx.Principal](
	pService principal.Service[T],
	apiSecretKey string,
	handler goservectx.ApiHandler[T],
) Service[T] {
	return &BaseService[T]{
		Service:         encryptor.New([]byte(apiSecretKey)),
		PService:        pService,
		ErrorHandler:    handler,
		claimsExtractor: defaultClaimsExtractor,
	}
}
