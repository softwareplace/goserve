package context

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	apiAccessContextKey = "apiAccessContext"
	XApiKey             = "X-Api-Key"
	Authorization       = "Authorization"
)

type Principal interface {
	GetId() string
	GetRoles() []string
	EncryptedPassword() string
}

type ApiHandler[T Principal] interface {
	Handler(ctx *Request[T], err error, source string)
}

type SampleContext[T Principal] struct {
	ApiKey              string                 // The API key extracted from the HTTP request header. This is used to identify and authenticate the client making the API request.
	ApiKeyId            string                 // The unique identifier associated with the API key used in the request. Helps in tracking and logging API key usage.
	Authorization       string                 // The bearer token or other authorization token extracted from the HTTP request header. Used to authenticate the user of the API.
	Principal           *T                     // The principal context containing user or session-specific data, representing the authenticated entity in the request.
	sessionId           string                 // A unique identifier for the current API session. Used to track the lifecycle of requests in a session.
	AuthorizationClaims map[string]interface{} // A set of claims derived from the authorization token, providing additional metadata about the requester (e.g., roles, permissions, expiration).
	ApiKeyClaims        map[string]interface{} // A set of claims derived from the API key, detailing metadata associated with the key (e.g., usage limits, allowed resources).
	AccessId            string                 // A unique identifier representing access to a specific resource or API, often used for auditing or tracking access patterns.
	PathValues          map[string]string      // A map of route variables extracted from the request URL. Useful for handling dynamic URL parameters in the API endpoints.
	Headers             map[string][]string    // Headers contains a mapping of header keys to their respective values from the incoming HTTP request.
	QueryValues         map[string][]string    // A map containing the query parameters from the request URL. Each key corresponds to a query parameter name, and the value is a slice of strings representing the values of that parameter. Useful for processing and validating query parameters in API endpoints.
}

func (ctx *Request[T]) GetSample() SampleContext[T] {
	return SampleContext[T]{
		ApiKey:              ctx.ApiKey,
		ApiKeyId:            ctx.ApiKeyId,
		Authorization:       ctx.Authorization,
		Principal:           ctx.Principal,
		sessionId:           ctx.sessionId,
		AuthorizationClaims: ctx.AuthorizationClaims,
		ApiKeyClaims:        ctx.ApiKeyClaims,
		AccessId:            ctx.AccessId,
		PathValues:          ctx.PathValues,
	}
}

type Request[T Principal] struct {
	Writer              *http.ResponseWriter
	Request             *http.Request
	ApiKey              string                 // The API key extracted from the HTTP request header. This is used to identify and authenticate the client making the API request.
	ApiKeyId            string                 // The unique identifier associated with the API key used in the request. Helps in tracking and logging API key usage.
	Authorization       string                 // The bearer token or other authorization token extracted from the HTTP request header. Used to authenticate the user of the API.
	Principal           *T                     // The principal context containing user or session-specific data, representing the authenticated entity in the request.
	sessionId           string                 // A unique identifier for the current API session. Used to track the lifecycle of requests in a session.
	AuthorizationClaims map[string]interface{} // A set of claims derived from the authorization token, providing additional metadata about the requester (e.g., roles, permissions, expiration).
	ApiKeyClaims        map[string]interface{} // A set of claims derived from the API key, detailing metadata associated with the key (e.g., usage limits, allowed resources).
	AccessId            string                 // A unique identifier representing access to a specific resource or API, often used for auditing or tracking access patterns.
	PathValues          map[string]string      // A map of route variables extracted from the request URL. Useful for handling dynamic URL parameters in the API endpoints.
	Headers             map[string][]string    // Headers contains a mapping of header keys to their respective values from the incoming HTTP request.
	QueryValues         map[string][]string    // A map containing the query parameters from the request URL. Each key corresponds to a query parameter name, and the value is a slice of strings representing the values of that parameter. Useful for processing and validating query parameters in API endpoints.
	Completed           bool                   // Completed indicates whether the task or process has been finished successfully or not.
}

// Of retrieves the Request object from the request's context if it already exists.
// If no such object exists, it creates a new instance of Request with the given writer, request,
// and reference for logging or tracing purposes.
//
// This function enhances the context of the current HTTP request with additional API-related information,
// such as API key, authorization token, and a unique session ID. The new context or the retrieved existing
// context is linked to the request to facilitate data sharing throughout the request's lifecycle.
//
// Type Parameters:
//   - T: A type that implements the Principal interface, which facilitates the storage
//     and management of additional API-related data for the request.
//
// Parameters:
//   - w: The http.ResponseWriter used to construct the response for the client.
//   - r: The *http.Request representing the HTTP request from the client.
//   - reference: A string value for logging or reference purposes.
//
// Returns:
//   - A pointer to the Request of type T, which contains relevant API-related data.
//
// Example usage:
//
//	ctx := Of[MyContextData](w, r, "MyReference")
//	ctx.GetSessionId() // Access session id
func Of[T Principal](w http.ResponseWriter, r *http.Request, reference string) *Request[T] {
	currentContext := r.Context().Value(apiAccessContextKey)

	if currentContext != nil {
		ctx := currentContext.(*Request[T])
		ctx.updateContext(r)
		return ctx
	}

	return createNewContext[T](w, r, reference)
}

// Flush clears all the fields in the Request, effectively resetting
// the context to its default state. This can be useful to prevent accidental
// reuse of sensitive data or to prepare for cleanup at the end of a request.
//
// This function clears sensitive information such as API key, authorization
// tokens, claims, and other metadata. It also nils out the Writer and Request
// pointers to avoid accidental usage after the flush.
//
// Usage Example:
//
//	ctx := Of[MyPrincipalContext](w, r, "ExampleReference")
//	// Process request here...
//	ctx.Flush() // Reset the context to its default state for cleanup.
func (ctx *Request[T]) Flush() {
	apiRequestContext := context.WithValue(ctx.Request.Context(), apiAccessContextKey, nil)
	ctx.Request = ctx.Request.WithContext(apiRequestContext)

	ctx.Writer = nil
	ctx.Request = nil
	ctx.AuthorizationClaims = nil
	ctx.ApiKeyClaims = nil
	ctx.Authorization = ""
	ctx.ApiKeyId = ""
	ctx.ApiKey = ""
	ctx.sessionId = ""
	ctx.AccessId = ""
	ctx.Principal = nil
}

func createNewContext[T Principal](
	w http.ResponseWriter,
	r *http.Request, reference string,
) *Request[T] {
	w.Header().Set("Content-Type", "application/json")
	ctx := Request[T]{
		Writer:        &w,
		Request:       r,
		PathValues:    mux.Vars(r),
		QueryValues:   r.URL.Query(),
		Headers:       r.Header,
		sessionId:     uuid.New().String(),
		ApiKey:        r.Header.Get(XApiKey),
		Authorization: r.Header.Get(Authorization),
	}

	log.Printf("%s -> initialized a context with session id: %s", reference, ctx.sessionId)
	ctx.updateContext(r)
	return &ctx
}

func (ctx *Request[T]) updateContext(r *http.Request) {
	apiRequestContext := context.WithValue(ctx.Request.Context(), apiAccessContextKey, ctx)
	ctx.Request = r.WithContext(apiRequestContext)
}

// Next forwards the request to the next HTTP handler in the middleware chain.
// It ensures the current Request is preserved during the request processing.
func (ctx *Request[T]) Next(next http.Handler) {
	next.ServeHTTP(*ctx.Writer, ctx.Request)
}

func (ctx *Request[T]) Done() {
	ctx.Completed = true
}

func (ctx *Request[T]) Write(body any, status int) {
	if !ctx.Completed {
		(*ctx.Writer).WriteHeader(status)
		err := json.NewEncoder(*ctx.Writer).Encode(body)
		if err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		ctx.Done()
	}
}
