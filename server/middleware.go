package server

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/goserve/context"
	errorhandler "github.com/softwareplace/goserve/error"
	"net/http"
	"time"
)

func (a *baseServer[T]) RegisterMiddleware(middleware ApiMiddleware[T], name string) Api[T] {
	a.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := apicontext.Of[T](w, r, name)
			if middleware(ctx) {
				ctx.Next(next)
			}
		})
	})
	return a
}

// rootAppMiddleware logs each incoming request's method, path, and remote address
func rootAppMiddleware[T apicontext.Principal](next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx *apicontext.Request[T]

		errorhandler.Handler(func() {
			start := time.Now() // Record the start time
			ctx = apicontext.Of[T](w, r, "MIDDLEWARE/ROOT_APP")
			queryParam := ""
			if r.URL.RawQuery != "" {
				queryParam = "?" + r.URL.RawQuery
			}

			log.Printf("[%s]:: Incoming request: %s %s from %s", ctx.GetSessionId(), r.Method, r.URL.Path+queryParam, r.RemoteAddr)

			ctx.Next(next)

			duration := time.Since(start)

			log.Printf("[%s]:: => request processed: %s %s in %v",
				ctx.GetSessionId(),
				r.Method,
				r.URL.Path+queryParam,
				duration,
			)

		}, func(err error) {
			onError(err, w)
		})

		defer func() {
			errorhandler.Handler(ctx.Flush, func(err error) {
				log.Errorf("Error flushing context: %v", err)
			})
			ctx = nil
		}()
	})
}

func onError(err any, w http.ResponseWriter) {
	log.Errorf("Error processing request: %+v", err)

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
