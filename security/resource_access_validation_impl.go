package security

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/security/router"
)

type defaultResourceAccessHandler[T goservectx.Principal] struct {
	handler *goserveerror.ApiHandler[T]
}

func (a *defaultResourceAccessHandler[T]) HasResourceAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := goservectx.Of[T](w, r, goserveerror.SecurityValidatorResourceAccess)

		if router.IsPublicPath(ctx.Request.Method, ctx.Request.URL.Path) {
			ctx.Next(next)
			return
		}

		if a.HasResourceAccessRight(*ctx) {
			ctx.Next(next)
			return
		}

		if a.handler != nil {
			(*a.handler).Handler(ctx, nil, goserveerror.SecurityValidatorResourceAccess)
			return
		}

		ctx.Forbidden("Access denied")
	})
}

func (a *defaultResourceAccessHandler[T]) HasResourceAccessRight(ctx goservectx.Request[T]) bool {
	requiredRoles, isRoleRequired := router.GetRolesForPath(ctx.Request.Method, ctx.Request.URL.Path)
	userRoles := (*ctx.Principal).GetRoles()

	if userRoles == nil || len(userRoles) == 0 {
		log.Errorf("Error: User roles are nil. Required roles: %v", requiredRoles)
		return false
	}

	for _, requiredRole := range requiredRoles {
		for _, userRole := range userRoles {
			if requiredRole == userRole {
				return true
			}
		}
	}

	log.Errorf("Error: User roles are nil. Required roles: [%s] but found [%v]",
		strings.Join(requiredRoles, ", "),
		strings.Join(userRoles, ", "),
	)

	return !isRoleRequired
}
