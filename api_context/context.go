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

type ApiPrincipalContext interface {
	GetSalt() string
	GetRoles() []string
}

type ApiRequestContext[T ApiPrincipalContext] struct {
	Writer              *http.ResponseWriter
	Request             *http.Request
	ApiKey              string                 // The API key extracted from the HTTP request header. This is used to identify and authenticate the client making the API request.
	ApiKeyId            string                 // The unique identifier associated with the API key used in the request. Helps in tracking and logging API key usage.
	Authorization       string                 // The bearer token or other authorization token extracted from the HTTP request header. Used to authenticate the user of the API.
	Principal           T                      // The principal context containing user or session-specific data, representing the authenticated entity in the request.
	sessionId           string                 // A unique identifier for the current API session. Used to track the lifecycle of requests in a session.
	AuthorizationClaims map[string]interface{} // A set of claims derived from the authorization token, providing additional metadata about the requester (e.g., roles, permissions, expiration).
	ApiKeyClaims        map[string]interface{} // A set of claims derived from the API key, detailing metadata associated with the key (e.g., usage limits, allowed resources).
	AccessId            string                 // A unique identifier representing access to a specific resource or API, often used for auditing or tracking access patterns.
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
//   - T: A type that implements the ApiPrincipalContext interface, which facilitates the storage
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
func Of[T ApiPrincipalContext](w http.ResponseWriter, r *http.Request, reference string) *ApiRequestContext[T] {
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

func createNewContext[T ApiPrincipalContext](
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
