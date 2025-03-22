package server

import "github.com/softwareplace/http-utils/security"

func (a *apiRouterHandlerImpl[T]) ApiSecretAccessHandler(apiSecretAccessHandler security.ApiSecretAccessHandler[T]) Api[T] {
	a.apiSecretAccessHandler = apiSecretAccessHandler
	return a.RegisterMiddleware(apiSecretAccessHandler.HandlerSecretAccess, security.ApiSecretAccessHandlerName)
}

func (a *apiRouterHandlerImpl[T]) ApiSecurityService(apiSecurityService security.Service[T]) Api[T] {
	a.apiSecurityService = apiSecurityService
	return a.RegisterMiddleware(apiSecurityService.AuthorizationHandler, security.ApiSecurityHandlerName)
}
