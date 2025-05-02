package jwt

import (
	log "github.com/sirupsen/logrus"
	goservectx "github.com/softwareplace/goserve/context"
)

func (a *BaseService[T]) HandlerErrorOrElse(
	ctx *goservectx.Request[T],
	error error,
	executionContext string,
	handlerNotFound func(),
) {
	if a.ErrorHandler != nil {
		a.ErrorHandler.Handler(ctx, error, executionContext)
		return
	}

	if handlerNotFound != nil {
		handlerNotFound()
		return
	}

	log.Errorf("DEFAULT/ERROR/HANDLER:: Failed to handle the request. Error: %s", error.Error())
	ctx.InternalServerError("Failed to handle the request. Please try again.")
}
