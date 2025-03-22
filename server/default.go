package server

import (
	apicontext "github.com/softwareplace/http-utils/context"
)

func Default(topMiddlewares ...ApiMiddleware[*apicontext.DefaultContext]) ApiRouterHandler[*apicontext.DefaultContext] {
	return CreateApiRouter[*apicontext.DefaultContext](topMiddlewares...)
}
