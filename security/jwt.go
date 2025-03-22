package security

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"github.com/softwareplace/http-utils/apicontext"
	"net/http"
	"time"
)

const (
	JWTLoadPrincipalError = "JWT/LOAD_PRINCIPAL_ERROR"
	JWTExtractClaimsError = "JWT/EXTRACT_CLAIMS_ERROR"
)

type ApiJWTInfo struct {
	Client     string
	Key        string
	Expiration time.Duration
	Scopes     []string
	PublicKey  *string
}

func (a *apiSecurityServiceImpl[T]) Principal(
	ctx *apicontext.ApiRequestContext[T],
) bool {
	success := a.PService.LoadPrincipal(ctx)

	if !success {
		a.handlerErrorOrElse(ctx, nil, JWTLoadPrincipalError, func() {
			ctx.Error("AuthorizationHandler failed", http.StatusForbidden)
		})

		return success
	}

	return success
}

func (a *apiSecurityServiceImpl[T]) ExtractJWTClaims(ctx *apicontext.ApiRequestContext[T]) bool {

	token, err := jwt.Parse(ctx.Authorization, func(token *jwt.Token) (interface{}, error) {
		return a.Secret(), nil
	})

	if err != nil {
		log.Printf("JWT/PARSE: AuthorizationHandler failed: %v", err)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.AuthorizationClaims = claims

		requester, err := a.Decrypt(claims["request"].(string))

		if err != nil {
			log.Printf("%s: AuthorizationHandler failed: %v", JWTExtractClaimsError, err)
			a.handlerErrorOrElse(ctx, err, JWTExtractClaimsError, func() {
				ctx.Error("AuthorizationHandler failed", http.StatusForbidden)
			})
			return false
		}

		ctx.AccessId = requester

		return true
	}

	log.Printf("JWT/CLAIMS_EXTRACT: failed with error_handler: %v", err)

	a.handlerErrorOrElse(ctx, err, JWTExtractClaimsError, func() {
		ctx.Error("AuthorizationHandler failed", http.StatusForbidden)
	})

	return false
}

func (a *apiSecurityServiceImpl[T]) JWTClaims(ctx *apicontext.ApiRequestContext[T]) (map[string]interface{}, error) {
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

func (a *apiSecurityServiceImpl[T]) GenerateJWT(data T, duration time.Duration) (*JwtResponse, error) {
	expiration := time.Now().Add(duration).Unix()
	requestBy, err := a.Encrypt(data.RequesterId())

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

	return &JwtResponse{
		Token:   signedToken,
		Expires: int(expiration),
	}, err
}
