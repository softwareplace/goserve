package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/security/encryptor"
	"github.com/softwareplace/goserve/utils"
	"time"
)

func (a *impl[T]) ExtractJWTClaims(ctx *goservectx.Request[T]) bool {
	token, err := a.Parse(ctx.Authorization)

	if err != nil {
		log.Errorf("JWT/PARSE: AuthorizationHandler failed: %+v", err)
		a.HandlerErrorOrElse(ctx, err, goserveerror.ExtractClaimsError, nil)
		return false
	}

	if claims, ok := a.Get(token); ok {
		ctx.AuthorizationClaims = claims

		isJwtClaimsEncryptionEnabled := encryptor.JwtClaimsEncryptionEnabled()

		requester := claims[SUB].(string)

		if isJwtClaimsEncryptionEnabled {
			requester, err = a.Decrypt(requester)
		}

		if err != nil {
			log.Errorf("%s: AuthorizationHandler failed: %+v", goserveerror.ExtractClaimsError, err)
			a.HandlerErrorOrElse(ctx, err, goserveerror.ExtractClaimsError, nil)
			return false
		}

		ctx.AccessId = requester

		return true
	}

	log.Errorf("JWT/CLAIMS_EXTRACT: failed with error: %+v", err)

	a.HandlerErrorOrElse(ctx, err, goserveerror.ExtractClaimsError, nil)
	return false
}

func (a *impl[T]) Generate(data T, duration time.Duration) (*Response, error) {
	return a.From(data.GetId(), data.GetRoles(), duration)
}

func (a *impl[T]) From(sub string, roles []string, duration time.Duration) (*Response, error) {
	if sub == "" {
		return nil, fmt.Errorf("sub cannot be empty")
	}

	iat := time.Now()
	expiration := iat.Add(duration).Unix()

	isJwtClaimsEncryptionEnabled := encryptor.JwtClaimsEncryptionEnabled()

	var err error
	var requestBy string
	var claimRoles []string

	if isJwtClaimsEncryptionEnabled {
		requestBy, err = a.Encrypt(sub)
		if err != nil {
			return nil, err
		}

		for _, role := range roles {
			var encryptedRole string
			encryptedRole, err = a.Encrypt(role)
			if err != nil {
				return nil, err
			}
			claimRoles = append(claimRoles, encryptedRole)
		}
	} else {
		requestBy = sub
		claimRoles = roles
	}

	claims := a.Create(requestBy, claimRoles, expiration, iat, a.Issuer())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(a.Secret())

	return &Response{
		JWT:      signedToken,
		Expires:  int(expiration),
		IssuedAt: int(iat.Unix()),
	}, err
}

func (a *impl[T]) Issuer() string {
	return utils.GetEnvOrDefault("JWT_ISSUER", "")
}

func (a *impl[T]) Decode(tokenString string) (map[string]interface{}, error) {
	token, err := a.Parse(tokenString)

	if err != nil {
		return nil, err
	}

	if claims, ok := a.Get(token); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims structure")
}

func (a *impl[T]) Parse(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return a.Secret(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	return token, nil
}

func (a *impl[T]) HandlerErrorOrElse(
	ctx *goservectx.Request[T],
	error error,
	executionContext string,
	handlerNotFound func(),
) {
	if a.ErrorHandler != nil {
		a.ErrorHandler.Handler(ctx, error, executionContext)
		return
	}

	if handlerNotFound != nil {
		handlerNotFound()
		return
	}

	log.Errorf("DEFAULT/ERROR/HANDLER:: Failed to handle the request. Error: %s", error.Error())
	ctx.InternalServerError("Failed to handle the request. Please try again.")
}
