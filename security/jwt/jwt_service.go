package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security/encryptor"
	"github.com/softwareplace/goserve/security/jwt/claims"
	"github.com/softwareplace/goserve/security/jwt/internal/impl"
	"github.com/softwareplace/goserve/security/jwt/response"
	"github.com/softwareplace/goserve/security/jwt/validate"
)

type Service[T goservectx.Principal] interface {
	encryptor.Service
	claims.Claims
	validate.Validate

	// Generate creates a new JWT Response based on the provided user and duration.
	// It returns the generated Response or an error if the process fails.
	Generate(user T, duration time.Duration) (*response.Response, error)

	// From generating a Response containing a JWT for the given subject, roles
	// and duration or returns an error if it fails.
	From(sub string, roles []string, duration time.Duration) (*response.Response, error)

	// Issuer returns the identifier of the entity responsible for issuing the JWT tokens in the service.
	Issuer() string

	// Decode extracts claims from a given JWT token string.
	// It validates the token and parses its claims into a map[string]interface{}.
	// Returns an error if the token is invalid or if the claims' structure is incorrect.
	Decode(tokenString string) (map[string]interface{}, error)

	// Decrypted extracts claims from a given JWT token string and attempts to decrypt ISS, SUB and AUD values.
	// If JWT claims encryption is enabled, it will decrypt the encrypted ISS, SUB and AUD claims values.
	// If encryption is disabled, it will return the original values.
	//
	// Parameters:
	//   - jwt: The JWT token string to extract and decrypt claims from.
	//
	// Returns:
	//   - map[string]interface{}: The claims map with decrypted values for ISS,
	//     SUB and AUD if encryption is enabled
	//   - error: An error if token parsing or decryption fails
	Decrypted(jwt string) (map[string]interface{}, error)

	// Parse parses and validates a JWT token string.
	// It uses the secret key provided by the BaseService for token signing and validation.
	// Returns the parsed *jwt.Token or an error if the token cannot be parsed or is invalid.
	Parse(tokenString string) (*jwt.Token, error)

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
	ExtractJWTClaims(ctx *goservectx.Request[T]) bool

	HandlerErrorOrElse(
		ctx *goservectx.Request[T],
		error error,
		executionContext string,
		handlerNotFound func(),
	)
}

func NewClaims() claims.Claims {
	return impl.NewClaims()
}

func NewValidate(secret []byte) validate.Validate {
	return impl.NewValidate(secret)
}

func New[T goservectx.Principal](
	apiSecretKey string,
	handler goservectx.ApiHandler[T],
) Service[T] {
	secret := []byte(apiSecretKey)

	return impl.NewJwtServiceImpl(
		NewClaims(),
		NewValidate(secret),
		encryptor.New(secret),
		handler,
	)
}
