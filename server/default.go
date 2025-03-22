package server

import (
	"github.com/softwareplace/http-utils/apicontext"
)

func Default(topMiddlewares ...ApiMiddleware[*apicontext.DefaultContext]) ApiRouterHandler[*apicontext.DefaultContext] {
	return CreateApiRouter[*apicontext.DefaultContext](topMiddlewares...)
}
