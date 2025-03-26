package server

import (
	log "github.com/sirupsen/logrus"
	goservecontext "github.com/softwareplace/goserve/context"
	goserveerrohandler "github.com/softwareplace/goserve/error"
	"net/http"
)

const ErrorHandlerWrapper = "ERROR/HANDLER/WRAPPER"

func (a *baseServer[T]) errorHandlerWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := goservecontext.Of[T](w, r, ErrorHandlerWrapper)
		goserveerrohandler.Handler(func() {
			ctx.Next(next)
		}, func(err error) {
			a.onError(err, ctx)
		})
	})
}

func (a *baseServer[T]) onError(err error, ctx *goservecontext.Request[T]) {
	if a.errorHandler == nil {
		log.Errorf("Error processing request: %+v", err)
		ctx.Error("Failed to handle the request", http.StatusInternalServerError)
		return
	}

	a.errorHandler.Handler(ctx, err, ErrorHandlerWrapper)
}
