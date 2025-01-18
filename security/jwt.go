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

const (
	JWTLoadPrincipalError = "JWT/LOAD_PRINCIPAL_ERROR"
	JWTExtractClaimsError = "JWT/EXTRACT_CLAIMS_ERROR"
)

func (a *apiSecurityServiceImpl[T]) Validation(ctx api_context.ApiRequestContext[T], loadPrincipal func(ctx api_context.ApiRequestContext[T]) (T, bool)) (*T, bool) {
	success := a.ExtractJWTClaims(ctx)

	if !success {
		a.handlerErrorOrElse(&ctx, nil, JWTLoadPrincipalError, func() {
			ctx.Error("Authorization failed", http.StatusForbidden)
		})

		return nil, success
	}

	principal, success := loadPrincipal(ctx)

	if !success {
		a.handlerErrorOrElse(&ctx, nil, JWTLoadPrincipalError, func() {
			ctx.Error("Authorization failed", http.StatusForbidden)
		})

		return nil, success
	}

	(*a.PService).SetData(principal)
	return &principal, success
}

func (a *apiSecurityServiceImpl[T]) ExtractJWTClaims(ctx api_context.ApiRequestContext[T]) bool {

	token, err := jwt.Parse(ctx.Authorization, func(token *jwt.Token) (interface{}, error) {
		return a.Secret(), nil
	})

	if err != nil {
		log.Printf("JWT/PARSE: Authorization failed: %v", err)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		(*a.PService).SetAuthorizationClaims(claims)

		requester, err := a.Decrypt(claims["request"].(string))

		if err != nil {
			log.Printf("%s: Authorization failed: %v", JWTExtractClaimsError, err)
			a.handlerErrorOrElse(&ctx, err, JWTExtractClaimsError, func() {
				ctx.Error("Authorization failed", http.StatusForbidden)
			})
			return false
		}

		(*a.PService).SetAccessId(requester)

		return true
	}

	log.Printf("JWT/CLAIMS_EXTRACT: failed with error_handler: %v", err)

	a.handlerErrorOrElse(&ctx, err, JWTExtractClaimsError, func() {
		ctx.Error("Authorization failed", http.StatusForbidden)
	})

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

func (a *apiSecurityServiceImpl[T]) GenerateJWT(data T, duration time.Duration) (map[string]interface{}, error) {
	expiration := time.Now().Add(duration).Unix()
	requestBy, err := a.Encrypt(data.GetSalt())

	var encryptedRoles []string
	for _, role := range data.GetRoles() {
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
