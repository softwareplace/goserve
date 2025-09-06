package server

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	utils "github.com/softwareplace/goserve/utils"
)

var (
	healthResourceEnable = utils.GetBoolEnvOrDefault("HEALTH_RESOURCE_ENABLE", true)
	healthResourcePath   = utils.GetEnvOrDefault("HEALTH_RESOURCE_PATH", utils.APIContextPath()+"health")
)

func (a *baseServer[T]) RegisterMiddleware(middleware ApiMiddleware[T], name string) Api[T] {
	a.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := goservectx.Of[T](w, r, name)
			if middleware(ctx) {
				ctx.Next(next)
			}
		})
	})
	return a
}

// rootAppMiddleware logs each incoming request's method, path, and remote address
func rootAppMiddleware[T goservectx.Principal](next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if healthResourceEnable && r.URL.Path == healthResourcePath {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
			return
		}

		var ctx *goservectx.Request[T]

		goserveerror.Handler(func() {
			start := time.Now() // Record the start time
			ctx = goservectx.Of[T](w, r, "MIDDLEWARE/ROOT_APP")

			uri := r.URL.RequestURI()
			log.Printf("[%s]:: Incoming request: %s %s from %s", ctx.GetSessionId(), r.Method, uri, r.RemoteAddr)

			ctx.Next(next)

			duration := time.Since(start)

			log.Printf("[%s]:: => request processed: %s %s in %v",
				ctx.GetSessionId(),
				r.Method,
				uri,
				duration,
			)

		}, func(err error) {
			onError(err, w)
		})

		defer func() {
			goserveerror.Handler(ctx.Flush, func(err error) {
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
