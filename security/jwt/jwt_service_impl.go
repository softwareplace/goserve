package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/utils"
	"net/http"
	"time"
)

func (a *serviceImpl[T]) Principal(
	ctx *apicontext.Request[T],
) bool {
	success := a.PService.LoadPrincipal(ctx)

	if !success {
		a.HandlerErrorOrElse(ctx, nil, LoadPrincipalError, func() {
			ctx.Error("AuthorizationHandler failed", http.StatusForbidden)
		})

		return success
	}

	return success
}

func (a *serviceImpl[T]) ExtractJWTClaims(ctx *apicontext.Request[T]) bool {

	token, err := jwt.Parse(ctx.Authorization, func(token *jwt.Token) (interface{}, error) {
		return a.Secret(), nil
	})

	if err != nil {
		log.Errorf("JWT/PARSE: AuthorizationHandler failed: %+v", err)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.AuthorizationClaims = claims

		requester, err := a.Decrypt(claims[SUB].(string))

		if err != nil {
			log.Errorf("%s: AuthorizationHandler failed: %+v", ExtractClaimsError, err)
			a.HandlerErrorOrElse(ctx, err, ExtractClaimsError, func() {
				ctx.Error("AuthorizationHandler failed", http.StatusForbidden)
			})
			return false
		}

		ctx.AccessId = requester

		return true
	}

	log.Errorf("JWT/CLAIMS_EXTRACT: failed with error: %+v", err)

	a.HandlerErrorOrElse(ctx, err, ExtractClaimsError, func() {
		ctx.Error("AuthorizationHandler failed", http.StatusForbidden)
	})

	return false
}

func (a *serviceImpl[T]) JWTClaims(ctx *apicontext.Request[T]) (map[string]interface{}, error) {
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

func (a *serviceImpl[T]) GenerateJWT(data T, duration time.Duration) (*Response, error) {
	now := time.Now()
	expiration := now.Add(duration).Unix()
	requestBy, err := a.Encrypt(data.GetId())

	var encryptedRoles []string

	for _, role := range data.GetRoles() {
		encryptedRole, err := a.Encrypt(role)
		if err != nil {
			return nil, err
		}
		encryptedRoles = append(encryptedRoles, encryptedRole)
	}

	claims := jwt.MapClaims{
		SUB: requestBy,
		AUD: encryptedRoles,
		EXP: expiration,
		IAT: now.Unix(),
	}

	if issuer := a.Issuer(); issuer != "" {
		claims[ISS] = issuer
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(a.Secret())

	return &Response{
		JWT:      signedToken,
		Expires:  int(expiration),
		IssuedAt: int(now.Unix()),
	}, err
}

func (a *serviceImpl[T]) Issuer() string {
	return utils.GetEnvOrDefault("JWT_ISSUER", "")
}
