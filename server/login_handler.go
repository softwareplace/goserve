package server

import (
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/http-utils/context"
	errorhandler "github.com/softwareplace/http-utils/error"
	"github.com/softwareplace/http-utils/login"
	"github.com/softwareplace/http-utils/request"
)

func (a *baseServer[T]) Login(ctx *apicontext.Request[T]) {
	request.GetRequestBody(ctx, login.User{}, a.loginDataHandler, request.FailedToLoadBody[T])
}

func (a *baseServer[T]) ApiKeyGenerator(ctx *apicontext.Request[T]) {
	request.GetRequestBody(ctx, ApiKeyEntryData{}, a.apiKeyGeneratorDataHandler, request.FailedToLoadBody[T])
}

func (a *baseServer[T]) loginDataHandler(ctx *apicontext.Request[T], user login.User) {

	errorhandler.Handler(func() {
		decrypt, err := a.securityService.Decrypt(user.Password)
		if err != nil {
			log.Printf("LOGIN/DECRYPT: Failed to decrypt encryptor: %v", err)
		} else {
			user.Password = decrypt
		}

		principal, err := a.loginService.Login(user)

		if err != nil {
			log.Printf("LOGIN/LOGIN: Failed to login: %v", err)
			ctx.Forbidden("Login failed: Invalid username or password")
			return
		}

		isValidPassword := a.loginService.IsValidPassword(user, principal)

		if !isValidPassword {
			log.Printf("LOGIN/LOGIN: Failed to login: %v", err)
			ctx.Forbidden("Login failed: Invalid username or password")
			return
		}

		jwt, err := a.securityService.GenerateJWT(principal, a.loginService.TokenDuration())

		if err != nil {
			log.Printf("LOGIN/JWT: Failed to generate JWT: %v", err)
			ctx.InternalServerError("Login failed with internal server error. Please try again later.")
			return
		}

		ctx.Ok(jwt)
	}, func(err error) {
		log.Printf("LOGIN/HANDLER: Failed to handle request: %v", err)
		ctx.Forbidden("Login failed: Invalid username or password")
	})
}
