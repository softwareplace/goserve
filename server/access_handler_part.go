package server

import (
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/security/secret"
)

func (a *apiRouterHandlerImpl[T]) SecretService(service secret.Service[T]) Api[T] {
	a.secretService = service
	return a.RegisterMiddleware(service.HandlerSecretAccess, secret.AccessHandlerName)
}

func (a *apiRouterHandlerImpl[T]) SecurityService(service security.Service[T]) Api[T] {
	a.securityService = service
	return a.RegisterMiddleware(service.AuthorizationHandler, security.ApiSecurityHandlerName)
}
