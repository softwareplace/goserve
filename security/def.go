package security

import (
	"github.com/softwareplace/http-utils/api_context"
	"time"
)

type ApiSecurityService[T api_context.ApiContextData] interface {
	Secret() []byte
	GenerateApiSecretJWT(jwtInfo ApiJWTInfo) (string, error)
	ExtractJWTClaims(requestContext api_context.ApiRequestContext[T]) bool
	JWTClaims(ctx api_context.ApiRequestContext[T]) (map[string]interface{}, error)
	GenerateJWT(user T) (map[string]interface{}, error)
	Encrypt(key string) (string, error)
	Decrypt(encrypted string) (string, error)
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
