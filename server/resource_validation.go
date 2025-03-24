package server

import (
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security/principal"
	"net/http"
	"strings"
)

const (
	SecurityValidatorResourceAccess = "SECURITY/VALIDATOR/RESOURCE_ACCESS"
)

func (a *baseServer[T]) hasResourceAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := apicontext.Of[T](w, r, SecurityValidatorResourceAccess)

		if principal.IsPublicPath[T](*ctx) {
			ctx.Next(next)
			return
		}

		if a.hasResourceAccessRight(*ctx) {
			ctx.Next(next)
			return
		}

		if a.errorHandler != nil {
			a.errorHandler.Handler(ctx, nil, SecurityValidatorResourceAccess)
		}
	})
}

// hasResourceAccessRight checks if the user has the necessary roles to access the requested resource.
// It compares the roles assigned to the user with those required for the resource's path.
// If the path does not require any roles, the function returns true.
//
// Parameters:
//
//	ctx - The API request context containing user roles and request metadata.
//
// Returns:
//
//	bool - True if the user has the required roles or if the path does not require roles, false otherwise.
func (a *baseServer[T]) hasResourceAccessRight(ctx apicontext.Request[T]) bool {
	requiredRoles, isRoleRequired := principal.GetRolesForPath(ctx)
	userRoles := (*ctx.Principal).GetRoles()

	if userRoles == nil || len(userRoles) == 0 {
		log.Printf("Error: User roles are nil. Required roles: %v", requiredRoles)
		return false
	}

	for _, requiredRole := range requiredRoles {
		for _, userRole := range userRoles {
			if requiredRole == userRole {
				return true
			}
		}
	}

	log.Printf("Error: User roles are nil. Required roles: [%s] but found [%v]",
		strings.Join(requiredRoles, ", "),
		strings.Join(userRoles, ", "),
	)

	return !isRoleRequired
}
