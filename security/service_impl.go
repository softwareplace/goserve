package security

import (
	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security/router"
)

func (a *impl[T]) AuthorizationHandler(ctx *goservectx.Request[T]) (doNext bool) {
	if router.IsPublicPath(ctx.Request.Method, ctx.Request.URL.Path) {
		return true
	}

	if !a.ExtractJWTClaims(ctx) {
		ctx.Forbidden("Invalid JWT token")
		return false
	}

	return a.PService.LoadPrincipal(ctx)
}
