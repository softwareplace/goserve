package server

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/security/validator"
	"net/http"
)

func HasResourceAccess[T api_context.ApiContextData](next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := api_context.Of[T](w, r, "SECURITY/VALIDATOR/RESOURCE_ACCESS")

		if validator.IsPublicPath[T](*ctx) {
			ctx.Next(next)
			return
		}

		if validator.HasResourceAccessRight[T](*ctx) {
			ctx.Next(next)
			return
		}

		ctx.Error("Access denied", http.StatusForbidden)
	})
}
