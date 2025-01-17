package server

import (
	"log"
	"net/http"
	"os"
)

var (
	ContextPath = apiContextPath()
	Port        = apiPort()
)

func apiContextPath() string {
	if contextPath := os.Getenv("CONTEXT_PATH"); contextPath != "" {
		return contextPath
	}
	return "/server/app/v1/"
}
func apiPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}

func (a *apiRouterHandlerImpl[T]) StartServer() {
	log.Printf("Server started at http://localhost:%s%s", Port, ContextPath)
	log.Fatal(http.ListenAndServe(":"+Port, a.router))
}
