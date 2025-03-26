package server

import (
	log "github.com/sirupsen/logrus"
	goservecontext "github.com/softwareplace/goserve/context"
	goserveerrohandler "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/request"
	"github.com/softwareplace/goserve/security/login"
)

func (a *baseServer[T]) Login(ctx *goservecontext.Request[T]) {
	request.GetRequestBody(ctx, login.User{}, a.loginDataHandler, request.FailedToLoadBody[T])
}

func (a *baseServer[T]) ApiKeyGenerator(ctx *goservecontext.Request[T]) {
	request.GetRequestBody(ctx, ApiKeyEntryData{}, a.apiKeyGeneratorDataHandler, request.FailedToLoadBody[T])
}

func (a *baseServer[T]) loginDataHandler(ctx *goservecontext.Request[T], user login.User) {

	goserveerrohandler.Handler(func() {
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
