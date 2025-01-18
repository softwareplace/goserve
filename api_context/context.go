package api_context

import (
	"context"
	"github.com/google/uuid"
	"log"
	"net/http"
)

const (
	apiAccessContextKey = "apiAccessContext"
	XApiKey             = "X-Api-Key"
	Authorization       = "Authorization"
)

type ApiContextData interface {
	// SetAuthorizationClaims sets the authorization claims for the current context.
	// These claims typically contain information related to user permissions or access rights.
	//
	// Parameters:
	//   - authorizationClaims: A map containing key-value pairs of the authorization claims.
	SetAuthorizationClaims(authorizationClaims map[string]interface{})

	// SetApiKeyClaims sets the claims associated with the API key.
	// These claims may define the permissions or metadata associated with the API key.
	//
	// Parameters:
	//   - authorizationClaims: A map containing key-value pairs of the API key claims.
	SetApiKeyClaims(authorizationClaims map[string]interface{})

	// SetApiKeyId sets the unique identifier of the API key being used in the request.
	//
	// Parameters:
	//   - apiKeyId: A string representing the API key ID.
	SetApiKeyId(apiKeyId string)

	// SetAccessId sets the access ID associated with the current context.
	// The access ID typically represents a user or system entity allowed access to certain resources.
	//
	// Parameters:
	//   - accessId: A string representing the access ID.
	SetAccessId(accessId string)

	// Data adds or updates the contextual data for the current instance.
	// This method is useful for storing additional metadata or API-related information.
	//
	// Parameters:
	//   - data: An instance implementing the ApiContextData interface, containing the information to store.
	Data(data ApiContextData)

	// SetRoles defines the roles associated with the current context.
	// These roles are typically used for role-based access control in the system.
	//
	// Parameters:
	//   - roles: A slice of strings representing the roles.
	SetRoles(roles []string)

	// Salt retrieves a string salt value used for securing operations, such as hashing.
	//
	// Returns:
	//   - A string representing the salt value.
	Salt() string

	// Roles retrieves the list of roles associated with the current context.
	// These roles define the access rights or permissions of the current user or entity.
	//
	// Returns:
	//   - A slice of strings representing the roles.
	Roles() []string
}

type ApiRequestContext[T ApiContextData] struct {
	Writer        *http.ResponseWriter
	Request       *http.Request
	ApiKey        string
	Authorization string
	RequestData   T
	sessionId     string
}

// Of retrieves the ApiRequestContext object from the request's context if it already exists.
// If no such object exists, it creates a new instance of ApiRequestContext with the given writer, request,
// and reference for logging or tracing purposes.
//
// This function enhances the context of the current HTTP request with additional API-related information,
// such as API key, authorization token, and a unique session ID. The new context or the retrieved existing
// context is linked to the request to facilitate data sharing throughout the request's lifecycle.
//
// Type Parameters:
//   - T: A type that implements the ApiContextData interface, which facilitates the storage
//     and management of additional API-related data for the request.
//
// Parameters:
//   - w: The http.ResponseWriter used to construct the response for the client.
//   - r: The *http.Request representing the HTTP request from the client.
//   - reference: A string value for logging or reference purposes.
//
// Returns:
//   - A pointer to the ApiRequestContext of type T, which contains relevant API-related data.
//
// Example usage:
//
//	ctx := Of[MyContextData](w, r, "MyReference")
//	ctx.GetSessionId() // Access session id
func Of[T ApiContextData](w http.ResponseWriter, r *http.Request, reference string) *ApiRequestContext[T] {
	currentContext := r.Context().Value(apiAccessContextKey)

	if currentContext != nil {
		ctx := currentContext.(*ApiRequestContext[T])
		ctx.updateContext(r)
		return ctx
	}

	return createNewContext[T](w, r, reference)
}

func (ctx *ApiRequestContext[T]) Flush() {
	ctx.Writer = nil
	ctx.Request = nil
}

func createNewContext[T ApiContextData](
	w http.ResponseWriter,
	r *http.Request, reference string,
) *ApiRequestContext[T] {
	w.Header().Set("Content-Type", "application/json")
	ctx := ApiRequestContext[T]{
		Writer:        &w,
		Request:       r,
		sessionId:     uuid.New().String(),
		ApiKey:        r.Header.Get(XApiKey),
		Authorization: r.Header.Get(Authorization),
	}

	log.Printf("%s -> initialized a context with session id: %s", reference, ctx.sessionId)
	ctx.updateContext(r)
	return &ctx
}

func (ctx *ApiRequestContext[T]) updateContext(r *http.Request) {
	apiRequestContext := context.WithValue(ctx.Request.Context(), apiAccessContextKey, ctx)
	ctx.Request = r.WithContext(apiRequestContext)
}

func (ctx *ApiRequestContext[T]) GetSessionId() string {
	return ctx.sessionId
}

func (ctx *ApiRequestContext[T]) Next(next http.Handler) {
	next.ServeHTTP(*ctx.Writer, ctx.Request)
}
