package server

import (
	apicontext "github.com/softwareplace/http-utils/context"
	errorhandler "github.com/softwareplace/http-utils/error"
	"net/http"
)

const ErrorHandlerWrapper = "ERROR/HANDLER/WRAPPER"

func (a *apiRouterHandlerImpl[T]) errorHandlerWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := apicontext.Of[T](w, r, ErrorHandlerWrapper)
		errorhandler.Handler(func() {
			ctx.Next(next)
		}, func(err error) {
			if a.errorHandler != nil {
				a.errorHandler.Handler(ctx, err, ErrorHandlerWrapper)
			} else {
				ctx.Error("Failed to handle the request", http.StatusInternalServerError)
			}
		})
	})
}
