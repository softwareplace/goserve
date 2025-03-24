package jwt

import apicontext "github.com/softwareplace/goserve/context"

func (a *serviceImpl[T]) HandlerErrorOrElse(
	ctx *apicontext.Request[T],
	error error,
	executionContext string,
	handlerNotFound func(),
) {
	if a.ErrorHandler != nil {
		a.ErrorHandler.Handler(ctx, error, executionContext)
	} else {
		handlerNotFound()
	}
}
