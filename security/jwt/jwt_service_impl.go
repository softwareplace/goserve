package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security/encryptor"
	"github.com/softwareplace/goserve/utils"
	"net/http"
	"time"
)

func (a *BaseService[T]) Principal(
	ctx *goservectx.Request[T],
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

func (a *BaseService[T]) ExtractJWTClaims(ctx *goservectx.Request[T]) bool {
	token, err := a.Parse(ctx.Authorization)

	if err != nil {
		log.Errorf("JWT/PARSE: AuthorizationHandler failed: %+v", err)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.AuthorizationClaims = claims

		isJwtClaimsEncryptionEnabled := encryptor.JwtClaimsEncryptionEnabled()

		var err error
		requester := claims[SUB].(string)

		if isJwtClaimsEncryptionEnabled {
			requester, err = a.Decrypt(claims[SUB].(string))
		}

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

func (a *BaseService[T]) Generate(data T, duration time.Duration) (*Response, error) {
	return a.From(data.GetId(), data.GetRoles(), duration)
}

func (a *BaseService[T]) From(sub string, roles []string, duration time.Duration) (*Response, error) {
	now := time.Now()
	expiration := now.Add(duration).Unix()

	isJwtClaimsEncryptionEnabled := encryptor.JwtClaimsEncryptionEnabled()

	var err error
	var requestBy string
	var encryptedRoles []string

	if isJwtClaimsEncryptionEnabled {
		requestBy, err = a.Encrypt(sub)
		if err != nil {
			return nil, err
		}

		for _, role := range roles {
			encryptedRole, err := a.Encrypt(role)
			if err != nil {
				return nil, err
			}
			encryptedRoles = append(encryptedRoles, encryptedRole)
		}
	} else {
		requestBy = sub
		encryptedRoles = roles
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

func (a *BaseService[T]) Issuer() string {
	return utils.GetEnvOrDefault("JWT_ISSUER", "")
}

func (a *BaseService[T]) Decode(tokenString string) (map[string]interface{}, error) {
	token, err := a.Parse(tokenString)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims structure")
}

func (a *BaseService[T]) Parse(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return a.Secret(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	return token, nil
}
