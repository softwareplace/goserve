package server

import (
	log "github.com/sirupsen/logrus"
	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	"net/http"
)

func (a *baseServer[T]) errorHandlerWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := goservectx.Of[T](w, r, goserveerror.HandlerWrapper)
		goserveerror.Handler(func() {
			ctx.Next(next)
		}, func(err error) {
			a.onError(err, ctx)
		})
	})
}

func (a *baseServer[T]) onError(err error, ctx *goservectx.Request[T]) {
	if a.errorHandler == nil {
		log.Errorf("Error processing request: %+v", err)
		ctx.Error("Failed to handle the request", http.StatusInternalServerError)
		return
	}

	a.errorHandler.Handler(ctx, err, goserveerror.HandlerWrapper)
}
