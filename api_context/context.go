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
	SetAuthorizationClaims(authorizationClaims map[string]interface{})
	SetApiKeyId(apiKeyId string)
	SetAccessId(accessId string)
	Data(data ApiContextData)
	SetRoles(roles []string)
	Salt() string
	Roles() []string
}

type ApiRequestContext[T ApiContextData] struct {
	Writer        http.ResponseWriter
	Request       *http.Request
	ApiKey        string
	Authorization string
	RequestData   T
	sessionId     string
}

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
		Writer:        w,
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
	next.ServeHTTP(ctx.Writer, ctx.Request)
}
