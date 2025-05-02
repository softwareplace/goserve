package server

import (
	log "github.com/sirupsen/logrus"
	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/request"
	"github.com/softwareplace/goserve/security/login"
	"github.com/softwareplace/goserve/security/model"
)

// Login handles user login requests by processing the request body and
// delegating to the loginDataHandler function. It ensures proper error
// handling in cases where the request body cannot be loaded.
func (a *baseServer[T]) Login(ctx *goservectx.Request[T]) {
	request.GetRequestBody(ctx, login.User{}, a.loginDataHandler, request.FailedToLoadBody[T])
}

// ApiKeyGenerator handles the generation of API keys by processing the request body
// and delegating to the secret service handler. It ensures proper error handling
// for cases where the request body cannot be loaded.
func (a *baseServer[T]) ApiKeyGenerator(ctx *goservectx.Request[T]) {
	request.GetRequestBody(ctx, model.ApiKeyEntryData{}, a.secretService.Handler, request.FailedToLoadBody[T])
}

func (a *baseServer[T]) loginDataHandler(ctx *goservectx.Request[T], user login.User) {

	goserveerror.Handler(func() {
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

		jwt, err := a.securityService.Generate(principal, a.loginService.TokenDuration())

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
