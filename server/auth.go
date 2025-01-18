package server

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/security"
	"log"
	"time"
)

type LoginEntryData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginService[T api_context.ApiPrincipalContext] interface {
	SecurityService() security.ApiSecurityService[T]
	Login(user LoginEntryData) (T, error)
	TokenDuration() time.Duration
}

func (a *apiRouterHandlerImpl[T]) Login(ctx *api_context.ApiRequestContext[T]) {
	GetRequestBody(ctx, LoginEntryData{}, a.loginDataHandler, FailedToLoadBody[T])
}

func (a *apiRouterHandlerImpl[T]) loginDataHandler(ctx *api_context.ApiRequestContext[T], loginEntryData LoginEntryData) {

	loginService := *a.loginService
	decrypt, err := loginService.SecurityService().Decrypt(loginEntryData.Password)
	if err != nil {
		log.Printf("LOGIN/DECRYPT: Failed to decrypt password: %v", err)
	} else {
		loginEntryData.Password = decrypt
	}

	login, err := loginService.Login(loginEntryData)

	if err != nil {
		log.Printf("LOGIN/LOGIN: Failed to login: %v", err)
		ctx.BadRequest("Login failed: Invalid username or password")
		return
	}

	jwt, err := loginService.SecurityService().GenerateJWT(login, loginService.TokenDuration())
	if err != nil {
		log.Printf("LOGIN/JWT: Failed to generate JWT: %v", err)
		ctx.InternalServerError("Login failed with internal server error. Please try again later.")
		return
	}

	ctx.Ok(jwt)
}
