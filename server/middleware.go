package server

import (
	"encoding/json"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/error_handler"
	"log"
	"net/http"
	"time"
)

// rootAppMiddleware logs each incoming request's method, path, and remote address
func rootAppMiddleware[T api_context.ApiContextData](next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		error_handler.Handler(func() {
			start := time.Now() // Record the start time
			ctx := api_context.Of[T](w, r, "MIDDLEWARE/ROOT_APP")

			log.Printf("[%s]:: Incoming request: %s %s from %s", ctx.GetSessionId(), r.Method, r.URL.Path, r.RemoteAddr)

			ctx.Next(next)

			duration := time.Since(start)
			log.Printf("[%s]:: => request processed: %s %s in %v",
				ctx.GetSessionId(),
				r.Method,
				r.URL.Path,
				duration,
			)

			error_handler.Handler(ctx.Flush, func(err error) {
				log.Printf("[%s]:: Error flushing response: %v", ctx.GetSessionId(), err)
			})
		}, func(err error) {
			onError(err, w)
		})
	})
}

func onError(err any, w http.ResponseWriter) {
	log.Printf("Error processing request: %+v", err)

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")

	responseBody := map[string]interface{}{
		"message":    "Failed to process request",
		"statusCode": http.StatusInternalServerError,
		"timestamp":  time.Now().UnixMilli(),
	}

	err = json.NewEncoder(w).Encode(responseBody)

	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
