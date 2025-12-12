package server

import (
	log "github.com/sirupsen/logrus"

	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/http"
	"github.com/softwareplace/goserve/security/login"
)

// Login handles user login requests by processing the request body and
// delegating to the loginDataHandler function. It ensures proper error
// handling in cases where the request body cannot be loaded.
func (a *baseServer[T]) Login(ctx *goservectx.Request[T]) {
	http.GetRequestBody(ctx, login.User{}, a.loginDataHandler, http.FailedToLoadBody[T])
}

func (a *baseServer[T]) loginDataHandler(ctx *goservectx.Request[T], user login.User) {

	goserveerror.Handler(func() {
		encryptedPass := false
		originalRequestPassword := user.Password

		decrypt, err := a.securityService.Decrypt(user.Password)
		if err != nil {
			log.Printf("LOGIN/DECRYPT: Failed to decrypt encryptor: %v", err)
		} else {
			encryptedPass = true
			user.Password = decrypt
		}

		principal, err := a.loginService.Login(user)

		if err != nil {
			log.Printf("LOGIN/LOGIN: Failed to login: %v", err)
			ctx.Forbidden("Login failed: Invalid username or password")
			return
		}

		isValidPassword := a.loginService.IsValidPassword(user, principal)

		if encryptedPass && !isValidPassword {
			// Retry validation with the original requested password
			user.Password = originalRequestPassword
			isValidPassword = a.loginService.IsValidPassword(user, principal)
		}

		if !isValidPassword {
			log.Printf("LOGIN/LOGIN: Failed to login: %v", err)
			ctx.Forbidden("Login failed: Invalid username or password")
			return
		}

		jwt, err := a.securityService.Generate(principal, a.loginService.TokenDuration(principal))

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
