package server

import "github.com/softwareplace/http-utils/security"

func (a *apiRouterHandlerImpl[T]) WithApiSecretAccessHandler(apiSecretAccessHandler security.ApiSecretAccessHandler[T]) ApiRouterHandler[T] {
	a.apiSecretAccessHandler = apiSecretAccessHandler
	return a.RegisterMiddleware(apiSecretAccessHandler.HandlerSecretAccess, security.ApiSecretAccessHandlerName)
}

func (a *apiRouterHandlerImpl[T]) WithApiSecurityService(apiSecurityService security.ApiSecurityService[T]) ApiRouterHandler[T] {
	a.apiSecurityService = apiSecurityService
	return a.RegisterMiddleware(apiSecurityService.AuthorizationHandler, security.ApiSecurityHandlerName)
}
