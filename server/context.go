package server

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

type ApiRequestContext struct {
	Writer        http.ResponseWriter
	Request       *http.Request
	ApiKey        string
	Authorization string
	RequestData   any
	sessionId     string
}

func Of(w http.ResponseWriter, r *http.Request, reference string) ApiRequestContext {
	currentContext := r.Context().Value(apiAccessContextKey)

	if currentContext != nil {
		var ctx *ApiRequestContext
		ctx = currentContext.(*ApiRequestContext)
		ctx.updateContext(r)
		return *ctx
	}

	return createNewContext(w, r, reference)
}

func (ctx *ApiRequestContext) Flush() {
	ctx.Writer = nil
	ctx.Request = nil
}

func createNewContext(w http.ResponseWriter, r *http.Request, reference string) ApiRequestContext {
	w.Header().Set("Content-Type", "application/json")
	ctx := ApiRequestContext{
		Writer:        w,
		Request:       r,
		sessionId:     uuid.New().String(),
		ApiKey:        r.Header.Get(XApiKey),
		Authorization: r.Header.Get(Authorization),
	}

	log.Printf("%s -> initialized a context with session id: %s", reference, ctx.sessionId)
	ctx.updateContext(r)
	return ctx
}

func (ctx *ApiRequestContext) updateContext(r *http.Request) {
	apiRequestContext := context.WithValue(ctx.Request.Context(), apiAccessContextKey, ctx)
	ctx.Request = r.WithContext(apiRequestContext)
}

func (ctx *ApiRequestContext) GetSessionId() string {
	return ctx.sessionId
}

func (ctx *ApiRequestContext) Next(next http.Handler) {
	next.ServeHTTP(ctx.Writer, ctx.Request)
}
