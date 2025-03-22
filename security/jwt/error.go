package jwt

import apicontext "github.com/softwareplace/http-utils/context"

func (a *serviceImpl[T]) handlerErrorOrElse(
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
