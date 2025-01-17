package server

import (
	"github.com/softwareplace/http-utils/api_context"
)

func Default(topMiddlewares ...ApiMiddleware[*api_context.DefaultContext]) ApiRouterHandler[*api_context.DefaultContext] {
	return New[*api_context.DefaultContext](topMiddlewares...)
}
