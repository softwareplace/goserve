package server

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func apiContextPath() string {
	if contextPath := os.Getenv("CONTEXT_PATH"); contextPath != "" {
		return contextPath
	}
	return "/"
}
func apiPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}

func (a *apiRouterHandlerImpl[T]) NotFoundHandler() ApiRouterHandler[T] {
	a.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("404 page not found: %s", r.URL.Path)

		swaggerPath := strings.TrimSuffix(a.contextPath, "/") + "/swagger"

		isSwaggerPath := strings.TrimSuffix(r.URL.Path, "/") == swaggerPath

		if a.swaggerIsEnabled && (r.URL.Path == a.contextPath || r.URL.Path == a.contextPath[:len(a.contextPath)-1] || isSwaggerPath) {
			a.goToSwaggerUi(w, r)
			return
		}

		log.Printf("Returning 404 page not found: %s", r.URL.Path)
		http.Error(w, "404 page not found", http.StatusNotFound)
	})
	return a
}

func (a *apiRouterHandlerImpl[T]) goToSwaggerUi(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, a.contextPath+"swagger/index.html", http.StatusMovedPermanently)
	log.Printf("Redirecting to swagger: %s", r.URL.Path)
}

func (a *apiRouterHandlerImpl[T]) CustomNotFoundHandler(handler func(w http.ResponseWriter, r *http.Request)) ApiRouterHandler[T] {
	a.router.NotFoundHandler = http.HandlerFunc(handler)
	return a
}

func (a *apiRouterHandlerImpl[T]) WithPort(port string) ApiRouterHandler[T] {
	a.port = port
	return a
}

func (a *apiRouterHandlerImpl[T]) WithContextPath(contextPath string) ApiRouterHandler[T] {
	a.contextPath = contextPath
	return a
}

func (a *apiRouterHandlerImpl[T]) StartServer() {
	if a.port == "" {
		a.port = apiPort()
	}

	if a.contextPath == "" {
		a.contextPath = apiContextPath()
	}

	log.Printf("Server started at http://localhost:%s%s", a.port, a.contextPath)
	log.Fatal(http.ListenAndServe(":"+a.port, a.router))
}
