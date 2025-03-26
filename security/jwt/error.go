package jwt

import goservecontext "github.com/softwareplace/goserve/context"

func (a *serviceImpl[T]) HandlerErrorOrElse(
	ctx *goservecontext.Request[T],
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
