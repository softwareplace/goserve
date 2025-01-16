package security

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/softwareplace/http-utils/api_context"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (a *apiSecurityServiceImpl[T]) Validation(
	ctx api_context.ApiRequestContext[T],
	next func(requestContext api_context.ApiRequestContext[T],
) (*T, bool)) (*T, bool) {
	success := a.ExtractJWTClaims(ctx)

	if !success {
		return nil, success
	}

	data, success := next(ctx)

	if !success {
		ctx.Error("Authorization failed", http.StatusForbidden)
		return nil, success
	}
	accessUserContext := ctx.RequestData

	accessUserContext.Data(*data)
	return data, success
}

func (a *apiSecurityServiceImpl[T]) ExtractJWTClaims(ctx api_context.ApiRequestContext[T]) bool {
	apiContext := ctx.RequestData

	token, err := jwt.Parse(ctx.Authorization, func(token *jwt.Token) (interface{}, error) {
		return a.Secret(), nil
	})

	if err != nil {
		log.Printf("JWT/PARSE: Authorization failed: %v", err)
		ctx.Error("Authorization failed", http.StatusForbidden)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		apiContext.SetAuthorizationClaims(claims)

		requester, err := a.Decrypt(claims["request"].(string))

		if err != nil {
			log.Printf("JWT/CLAIMS_EXTRACT: Authorization failed: %v", err)
			ctx.Error("Authorization failed", http.StatusForbidden)
			return false
		}

		apiContext.SetAccessId(requester)

		return true
	}

	log.Printf("JWT/CLAIMS_EXTRACT: failed with error_handler: %v", err)
	ctx.Error("Authorization failed", http.StatusForbidden)
	return false
}

func (a *apiSecurityServiceImpl[T]) JWTClaims(ctx api_context.ApiRequestContext[T]) (map[string]interface{}, error) {
	token, err := jwt.Parse(ctx.ApiKey, func(token *jwt.Token) (interface{}, error) {
		return a.Secret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("failed to extract jwt claims")
}

func (a *apiSecurityServiceImpl[T]) Secret() []byte {
	secret := a.ApiSecretAuthorization
	return []byte(secret)
}

// GenerateJWT creates a JWT token with the username and role
func (a *apiSecurityServiceImpl[T]) GenerateJWT(data T) (map[string]interface{}, error) {
	duration := time.Minute * 15
	expiration := time.Now().Add(duration).Unix()
	requestBy, err := a.Encrypt(data.Salt())

	var encryptedRoles []string
	for _, role := range data.Roles() {
		encryptedRole, err := a.Encrypt(role)
		if err != nil {
			return nil, err
		}
		encryptedRoles = append(encryptedRoles, encryptedRole)
	}

	claims := jwt.MapClaims{
		"request": requestBy,
		"scope":   encryptedRoles,
		"exp":     expiration,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(a.Secret())
	return map[string]interface{}{"token": signedToken, "expires": strconv.FormatInt(expiration, 10)}, err
}
