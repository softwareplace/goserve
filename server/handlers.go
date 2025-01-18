package server

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/error_handler"
	"net/http"
)

type ApiErrorHandler[T api_context.ApiContextData] interface {
	Handler(ctx *api_context.ApiRequestContext[T], err error, source string)
}

const ErrorHandlerWrapper = "ERROR/HANDLER/WRAPPER"

func (a *apiRouterHandlerImpl[T]) errorHandlerWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := api_context.Of[T](w, r, ErrorHandlerWrapper)
		error_handler.Handler(func() {
			ctx.Next(next)
		}, func(err error) {
			if a.errorHandler != nil {
				(*a.errorHandler).Handler(ctx, err, ErrorHandlerWrapper)
			} else {
				ctx.Error("Failed to handle the request", http.StatusInternalServerError)
			}
		})
	})
}
