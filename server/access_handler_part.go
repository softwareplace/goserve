package server

import (
	"github.com/softwareplace/goserve/security"
	"github.com/softwareplace/goserve/security/secret"
)

func (a *baseServer[T]) SecretService(service secret.Service[T]) Api[T] {
	a.secretService = service
	if a.apiSecretKeyGeneratorResourceEnable {
		scopes := service.RequiredScopes()
		a.Add(a.ApiKeyGenerator, "api-key/generate", "POST", scopes...)
	}
	return a.RegisterMiddleware(service.HandlerSecretAccess, secret.AccessHandlerName)
}

func (a *baseServer[T]) SecurityService(service security.Service[T]) Api[T] {
	a.securityService = service
	return a.RegisterMiddleware(service.AuthorizationHandler, security.ApiSecurityHandlerName).
		EmbeddedServer(func(Api[T]) { a.router.Use(service.HasResourceAccess) })
}
