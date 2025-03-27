package security

import (
	apicontext "github.com/softwareplace/goserve/context"
	"net/http"
)

type ResourceAccessValidation[T apicontext.Principal] interface {

	// HasResourceAccess wraps the given HTTP handler to enforce resource access control.
	//
	// This middleware checks if the request is accessing a public path or if the user has the required roles
	// for the requested resource. If access is granted, the request proceeds to the next handler. Otherwise,
	// it returns access denied response.
	//
	// Parameters:
	//   - next: The next http.Handler in the chain to call if access is granted.
	//
	// Returns:
	//   - http.Handler: A handler that wraps the provided handler with access control
	HasResourceAccess(next http.Handler) http.Handler

	// HasResourceAccessRight checks if the user has the necessary roles to access the requested resource.
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
	HasResourceAccessRight(ctx apicontext.Request[T]) bool
}
