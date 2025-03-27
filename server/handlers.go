package server

import (
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/goserve/context"
	errorhandler "github.com/softwareplace/goserve/error"
	"net/http"
)

func (a *baseServer[T]) errorHandlerWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := apicontext.Of[T](w, r, errorhandler.HandlerWrapper)
		errorhandler.Handler(func() {
			ctx.Next(next)
		}, func(err error) {
			a.onError(err, ctx)
		})
	})
}

func (a *baseServer[T]) onError(err error, ctx *apicontext.Request[T]) {
	if a.errorHandler == nil {
		log.Errorf("Error processing request: %+v", err)
		ctx.Error("Failed to handle the request", http.StatusInternalServerError)
		return
	}

	a.errorHandler.Handler(ctx, err, errorhandler.HandlerWrapper)
}
