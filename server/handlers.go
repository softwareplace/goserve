package server

import (
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/goserve/context"
	errorhandler "github.com/softwareplace/goserve/error"
	"net/http"
)

const ErrorHandlerWrapper = "ERROR/HANDLER/WRAPPER"

func (a *baseServer[T]) errorHandlerWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := apicontext.Of[T](w, r, ErrorHandlerWrapper)
		errorhandler.Handler(func() {
			ctx.Next(next)
		}, func(err error) {
			if a.errorHandler != nil {
				a.errorHandler.Handler(ctx, err, ErrorHandlerWrapper)
			} else {
				log.Errorf("Error processing request: %+v", err)
				ctx.Error("Failed to handle the request", http.StatusInternalServerError)
			}
		})
	})
}
