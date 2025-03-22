package context

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

// InternalServerError sends an HTTP 500 Internal Server Error response with a given message.
func (ctx *ApiRequestContext[T]) InternalServerError(message string) {
	ctx.Error(message, http.StatusInternalServerError)
}

// Forbidden sends an HTTP 403 Forbidden response with a given message.
func (ctx *ApiRequestContext[T]) Forbidden(message string) {
	ctx.Error(message, http.StatusForbidden)
}

// Unauthorized sends an HTTP 401 Unauthorized response with a standardized message.
func (ctx *ApiRequestContext[T]) Unauthorized() {
	ctx.Error("Unauthorized", http.StatusUnauthorized)
}

// InvalidInput sends an HTTP 400 Bad Request response with a standardized "Invalid input" message.
func (ctx *ApiRequestContext[T]) InvalidInput() {
	ctx.BadRequest("Invalid input")
}

// BadRequest sends an HTTP 400 Bad Request response with a given message.
func (ctx *ApiRequestContext[T]) BadRequest(message string) {
	ctx.Error(message, http.StatusBadRequest)
}

// Ok sends an HTTP 200 OK response with the provided body as the response payload.
func (ctx *ApiRequestContext[T]) Ok(body any) {
	ctx.Response(body, http.StatusOK)
}

// Created sends an HTTP 201 Created response with the provided body as the response payload.
func (ctx *ApiRequestContext[T]) Created(body any) {
	ctx.Response(body, http.StatusCreated)
}

// NoContent sends an HTTP 204 No Content response. The body is ignored as the status indicates no body.
func (ctx *ApiRequestContext[T]) NoContent(body any) {
	ctx.Response(body, http.StatusNoContent)
}

// NotFount sends an HTTP 404 Not Found response with the provided body as the response payload.
func (ctx *ApiRequestContext[T]) NotFount(body any) {
	ctx.Response(body, http.StatusNotFound)
}

// Error sends an HTTP error response with a status and a message. The response includes
// the timestamp and status code for debugging or informational purposes.
func (ctx *ApiRequestContext[T]) Error(message string, status int) {
	responseBody := map[string]interface{}{
		"message":    message,
		"statusCode": status,
		"timestamp":  time.Now().UnixMilli(),
	}

	ctx.Response(responseBody, status)
}

// Response sends a generic HTTP response with a given body and status code.
// It serializes the body to JSON and writes it to the response writer.
func (ctx *ApiRequestContext[T]) Response(body any, status int) {
	(*ctx.Writer).WriteHeader(status)
	err := json.NewEncoder(*ctx.Writer).Encode(body)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// WriteFile streams a file as a response for download, using a given file name.
// Sets appropriate headers for file attachment and streams the content.
func (ctx *ApiRequestContext[T]) WriteFile(file []byte, fileName string) error {
	writer := *ctx.Writer

	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	writer.Header().Set("Content-Type", "application/octet-stream")

	_, err := io.Copy(writer, bytes.NewReader(file))
	return err
}

// WriteReader streams a file-like reader as a response for download with a given file name.
// Sets appropriate headers for file attachment and streams the reader's content.
func (ctx *ApiRequestContext[T]) WriteReader(reader *bytes.Reader, fileName string) error {
	writer := *ctx.Writer

	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	writer.Header().Set("Content-Type", "application/octet-stream")

	_, err := io.Copy(writer, reader)
	return err
}
