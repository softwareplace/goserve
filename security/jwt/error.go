package jwt

import goservectx "github.com/softwareplace/goserve/context"

func (a *BaseService[T]) HandlerErrorOrElse(
	ctx *goservectx.Request[T],
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
