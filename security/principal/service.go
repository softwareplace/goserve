package principal

import (
	apicontext "github.com/softwareplace/http-utils/context"
)

type Service[T apicontext.Principal] interface {

	// LoadPrincipal loads the principal information for the given API request context.
	// This method is responsible for extracting and validating the user or session-specific
	// data from the incoming request, ensuring that the request is associated with a valid principal.
	//
	// Parameters:
	//   - ctx: The API request context containing the necessary metadata and headers.
	//
	// Returns:
	//   - A boolean value indicating whether the principal was successfully loaded.
	//	 Returns true if the principal is valid and loaded successfully; otherwise, false.
	LoadPrincipal(ctx *apicontext.Request[T]) bool
}
