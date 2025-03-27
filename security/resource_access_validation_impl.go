package security

import (
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/goserve/context"
	errorhandler "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/security/principal"
	"net/http"
	"strings"
)

type defaultResourceAccessHandler[T apicontext.Principal] struct {
	handler *errorhandler.ApiHandler[T]
}

func (a *defaultResourceAccessHandler[T]) HasResourceAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := apicontext.Of[T](w, r, errorhandler.SecurityValidatorResourceAccess)

		if principal.IsPublicPath[T](*ctx) {
			ctx.Next(next)
			return
		}

		if a.HasResourceAccessRight(*ctx) {
			ctx.Next(next)
			return
		}

		if a.handler != nil {
			(*a.handler).Handler(ctx, nil, errorhandler.SecurityValidatorResourceAccess)
			return
		}

		ctx.Forbidden("Access denied")
	})
}

func (a *defaultResourceAccessHandler[T]) HasResourceAccessRight(ctx apicontext.Request[T]) bool {
	requiredRoles, isRoleRequired := principal.GetRolesForPath(ctx)
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
