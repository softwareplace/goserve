package principal

import "github.com/softwareplace/http-utils/api_context"

type PService[T api_context.ApiPrincipalContext] interface {
	// GetRoles retrieves the roles assigned to the principal for the given API request context.
	//
	// Parameters:
	//   - ctx: The API request context containing user information and metadata.
	//
	// Returns:
	//   - A slice of strings representing the roles assigned to the principal.
	GetRoles(ctx api_context.ApiRequestContext[T]) []string

	// LoadPrincipal loads and returns the principal associated with the given API request context.
	//
	// Parameters:
	//   - ctx: The API request context containing authentication and user-related information.
	//
	// Returns:
	//   - A pointer to the principal of type T if successful, or nil otherwise.
	//   - A boolean indicating whether the operation succeeded.
	LoadPrincipal(ctx api_context.ApiRequestContext[T]) (T, bool)

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

	// SetData attaches context-specific data of type T to the current service instance.
	//
	// Parameters:
	//   - data: The context-specific data to associate with the principal.
	SetData(data T)
}
