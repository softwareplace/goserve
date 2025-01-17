package api_context

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (ctx *ApiRequestContext[T]) InternalServerError(message string) {
	ctx.Error(message, http.StatusInternalServerError)
}

func (ctx *ApiRequestContext[T]) Forbidden(message string) {
	ctx.Error(message, http.StatusForbidden)
}

func (ctx *ApiRequestContext[T]) Unauthorized() {
	ctx.Error("Unauthorized", http.StatusUnauthorized)
}

func (ctx *ApiRequestContext[T]) InvalidInput() {
	ctx.BadRequest("Invalid input")
}

func (ctx *ApiRequestContext[T]) BadRequest(message string) {
	ctx.Error(message, http.StatusBadRequest)
}

func (ctx *ApiRequestContext[T]) Ok(body any) {
	ctx.Response(body, http.StatusOK)
}

func (ctx *ApiRequestContext[T]) Created(body any) {
	ctx.Response(body, http.StatusCreated)
}

func (ctx *ApiRequestContext[T]) NoContent(body any) {
	ctx.Response(body, http.StatusNoContent)
}

func (ctx *ApiRequestContext[T]) NotFount(body any) {
	ctx.Response(body, http.StatusNotFound)
}

func (ctx *ApiRequestContext[T]) Error(message string, status int) {
	(*ctx.Writer).WriteHeader(status)
	responseBody := map[string]interface{}{
		"message":    message,
		"statusCode": status,
		"timestamp":  time.Now().UnixMilli(),
	}

	ctx.Response(responseBody, status)
}

func (ctx *ApiRequestContext[T]) Response(body any, status int) {
	(*ctx.Writer).WriteHeader(status)
	err := json.NewEncoder(*ctx.Writer).Encode(body)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
