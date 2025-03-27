package jwt

import (
	apicontext "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security/encryptor"
	"github.com/softwareplace/goserve/security/principal"
	"time"
)

const (
	IAT                = "iat"
	EXP                = "exp"
	AUD                = "aud"
	SUB                = "sub"
	ISS                = "iss"
	LoadPrincipalError = "JWT/LOAD_PRINCIPAL_ERROR"
	ExtractClaimsError = "JWT/EXTRACT_CLAIMS_ERROR"
)

type Response struct {
	JWT      string `json:"jwt"`
	Expires  int    `json:"expires"`
	IssuedAt int    `json:"issuedAt"`
}

type Service[T apicontext.Principal] interface {
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
	ExtractJWTClaims(requestContext *apicontext.Request[T]) bool

	// Generate creates a new JWT Response based on the provided user and duration.
	// It returns the generated Response or an error if the process fails.
	Generate(user T, duration time.Duration) (*Response, error)

	// From generates a Response containing a JWT for the given subject, roles, and duration or returns an error if it fails.
	From(sub string, roles []string, duration time.Duration) (*Response, error)

	// Issuer returns the identifier of the entity responsible for issuing the JWT tokens in the service.
	Issuer() string

	HandlerErrorOrElse(
		ctx *apicontext.Request[T],
		error error,
		executionContext string,
		handlerNotFound func(),
	)
}

type BaseService[T apicontext.Principal] struct {
	encryptor.Service
	PService     principal.Service[T]
	ErrorHandler apicontext.ApiHandler[T]
}

func New[T apicontext.Principal](
	pService principal.Service[T],
	apiSecretKey string,
	handler apicontext.ApiHandler[T],
) Service[T] {
	return &BaseService[T]{
		Service:      encryptor.New([]byte(apiSecretKey)),
		PService:     pService,
		ErrorHandler: handler,
	}
}
