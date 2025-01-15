package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (ctx *ApiRequestContext) InternalServerError(message string) {
	ctx.Error(message, http.StatusInternalServerError)
}

func (ctx *ApiRequestContext) Forbidden(message string) {
	ctx.Error(message, http.StatusForbidden)
}

func (ctx *ApiRequestContext) Unauthorized() {
	ctx.Error("Unauthorized", http.StatusUnauthorized)
}

func (ctx *ApiRequestContext) InvalidInput() {
	ctx.BadRequest("Invalid input")
}

func (ctx *ApiRequestContext) BadRequest(message string) {
	ctx.Error(message, http.StatusBadRequest)
}

func (ctx *ApiRequestContext) Ok(body any) {
	ctx.Response(body, http.StatusOK)
}

func (ctx *ApiRequestContext) Created(body any) {
	ctx.Response(body, http.StatusCreated)
}

func (ctx *ApiRequestContext) NoContent(body any) {
	ctx.Response(body, http.StatusNoContent)
}

func (ctx *ApiRequestContext) NotFount(body any) {
	ctx.Response(body, http.StatusNotFound)
}

func (ctx *ApiRequestContext) Error(message string, status int) {
	ctx.Writer.WriteHeader(status)
	responseBody := map[string]interface{}{
		"message":    message,
		"statusCode": status,
		"timestamp":  time.Now().UnixMilli(),
	}

	ctx.Response(responseBody, status)
}

func (ctx *ApiRequestContext) Response(body any, status int) {
	ctx.Writer.WriteHeader(status)
	err := json.NewEncoder(ctx.Writer).Encode(body)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
