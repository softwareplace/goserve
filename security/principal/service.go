package principal

import "github.com/softwareplace/http-utils/api_context"

type PService[T api_context.ApiPrincipalContext] interface {

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
	LoadPrincipal(ctx api_context.ApiRequestContext[T]) bool

	// SetAuthorizationClaims sets the authorization claims for the current service instance.
	// These claims are used to define access control rules and permissions.
	//
	// Parameters:
	//   - authorizationClaims: A map containing key-value pairs for authorization-related claims.
	SetAuthorizationClaims(authorizationClaims map[string]interface{})

	// SetApiKeyClaims configures the API key claims for the current service instance.
	// These claims are used to validate and authorize API key usage.
	//
	// Parameters:
	//   - authorizationClaims: A map containing key-value pairs for API key-related claims.
	SetApiKeyClaims(authorizationClaims map[string]interface{})

	// SetApiKeyId assigns an API key ID to the current service instance, used to track API key activities.
	//
	// Parameters:
	//   - apiKeyId: A string representing the API key ID.
	SetApiKeyId(apiKeyId string)

	// SetAccessId assigns an Access ID to the current service instance, used for resource and access tracking.
	//
	// Parameters:
	//   - accessId: A string representing the Access ID.
	SetAccessId(accessId string)
}
