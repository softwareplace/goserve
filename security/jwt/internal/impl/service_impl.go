package impl

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"

	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/env"
	goserveerror "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/security/encryptor"
	"github.com/softwareplace/goserve/security/jwt/claims"
	"github.com/softwareplace/goserve/security/jwt/constants"
	"github.com/softwareplace/goserve/security/jwt/response"
	"github.com/softwareplace/goserve/security/jwt/validate"
)

type jwtServiceImpl[T goservectx.Principal] struct {
	claims.Claims
	validate.Validate
	encryptor.Service
	ErrorHandler goservectx.ApiHandler[T]
}

func NewJwtServiceImpl[T goservectx.Principal](
	claims claims.Claims,
	validate validate.Validate,
	encryptorService encryptor.Service,
	errorHandler goservectx.ApiHandler[T],
) *jwtServiceImpl[T] {
	return &jwtServiceImpl[T]{
		Claims:       claims,
		Validate:     validate,
		Service:      encryptorService,
		ErrorHandler: errorHandler,
	}
}

func (a *jwtServiceImpl[T]) Decrypted(jwt string) (map[string]interface{}, error) {
	token, err := a.Decode(jwt)
	if err != nil {
		return nil, err
	}
	isJwtClaimsEncryptionEnabled := encryptor.JwtClaimsEncryptionEnabled()
	if isJwtClaimsEncryptionEnabled {

		if aud, containsKey := token[constants.AUD].([]interface{}); containsKey {
			var values []string
			for _, audV := range aud {
				decrypt, err := a.Decrypt(audV.(string))
				if err != nil {
					return nil, err
				}
				values = append(values, decrypt)
			}
			token[constants.AUD] = values
		}

		sub, err := a.DecryptClaimsValue(constants.SUB, token)

		if err != nil {
			return nil, err
		}

		token[constants.SUB] = sub
	}

	return token, nil
}

func (a *jwtServiceImpl[T]) DecryptClaimsValue(key string, claims map[string]interface{}) (interface{}, error) {
	value, containsKey := claims[key]

	if !containsKey {
		return nil, fmt.Errorf("key %s not found", key)
	}

	if isJwtClaimsEncryptionEnabled := encryptor.JwtClaimsEncryptionEnabled(); isJwtClaimsEncryptionEnabled {
		decrypt, err := a.Decrypt(value.(string))
		if err != nil {
			return nil, err
		}
		return decrypt, nil
	}
	return value, nil
}

func (a *jwtServiceImpl[T]) ExtractJWTClaims(ctx *goservectx.Request[T]) bool {
	token, err := a.Parse(ctx.Authorization)

	if err != nil {
		log.Errorf("JWT/PARSE: AuthorizationHandler failed: %+v", err)
		a.HandlerErrorOrElse(ctx, err, goserveerror.ExtractClaimsError, nil)
		return false
	}

	if claims, ok := a.Get(token); ok {
		ctx.AuthorizationClaims = claims

		isJwtClaimsEncryptionEnabled := encryptor.JwtClaimsEncryptionEnabled()

		requester := claims[constants.SUB].(string)

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

func (a *jwtServiceImpl[T]) Generate(data T, duration time.Duration) (*response.Response, error) {
	return a.From(data.GetId(), data.GetRoles(), duration)
}

func (a *jwtServiceImpl[T]) From(sub string, roles []string, duration time.Duration) (*response.Response, error) {
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

	return &response.Response{
		JWT:      signedToken,
		Expires:  int(expiration),
		IssuedAt: int(iat.Unix()),
	}, err
}

func (a *jwtServiceImpl[T]) Issuer() string {
	return env.GetEnvOrDefault("JWT_ISSUER", "")
}

func (a *jwtServiceImpl[T]) Decode(tokenString string) (map[string]interface{}, error) {
	token, err := a.Parse(tokenString)

	if err != nil {
		return nil, err
	}

	if claims, ok := a.Get(token); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims structure")
}

func (a *jwtServiceImpl[T]) Parse(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return a.Secret(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	return token, nil
}

func (a *jwtServiceImpl[T]) HandlerErrorOrElse(
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
