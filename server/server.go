package server

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/http-utils/context"
	errorhandler "github.com/softwareplace/http-utils/error"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/security/principal"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type apiRouterHandlerImpl[T apicontext.ApiPrincipalContext] struct {
	router                              *mux.Router
	principalService                    principal.PService[T]
	errorHandler                        errorhandler.ApiErrorHandler[T]
	loginService                        LoginService[T]
	apiSecurityService                  security.ApiSecurityService[T]
	apiSecretAccessHandler              security.ApiSecretAccessHandler[T]
	apiKeyGeneratorService              ApiKeyGeneratorService[T]
	server                              *http.Server // Add a server instance
	mu                                  sync.Mutex   // Add a mutex for thread safety
	swaggerIsEnabled                    bool
	loginResourceEnable                 bool
	apiSecretKeyGeneratorResourceEnable bool
	contextPath                         string
	port                                string
}

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

func (a *apiRouterHandlerImpl[T]) NotFoundHandler() Api[T] {
	a.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Errorf("404 page not found: %s", r.URL.Path)

		if r.Method == "GET" {
			swaggerPath := strings.TrimSuffix(a.contextPath, "/") + "/swagger"

			isSwaggerPath := strings.TrimSuffix(r.URL.Path, "/") == swaggerPath

			if a.swaggerIsEnabled && (r.URL.Path == a.contextPath || r.URL.Path == a.contextPath[:len(a.contextPath)-1] || isSwaggerPath) {
				a.goToSwaggerUi(w, r)
				return
			}
		}

		log.Warnf("Returning 404 page not found: %s", r.URL.Path)
		http.Error(w, "404 page not found", http.StatusNotFound)
	})
	return a
}

func (a *apiRouterHandlerImpl[T]) goToSwaggerUi(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, a.contextPath+"swagger/index.html", http.StatusMovedPermanently)
	log.Infof("Redirecting to swagger: %s", r.URL.Path)
}

func (a *apiRouterHandlerImpl[T]) CustomNotFoundHandler(handler func(w http.ResponseWriter, r *http.Request)) Api[T] {
	a.router.NotFoundHandler = http.HandlerFunc(handler)
	return a
}

func (a *apiRouterHandlerImpl[T]) Port(port string) Api[T] {
	a.port = port
	return a
}

func (a *apiRouterHandlerImpl[T]) ContextPath(contextPath string) Api[T] {
	a.contextPath = contextPath
	return a
}

func (a *apiRouterHandlerImpl[T]) StartServerInGoroutine() Api[T] {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.port == "" {
		a.port = apiPort()
	}

	if a.contextPath == "" {
		a.contextPath = apiContextPath()
	}

	addr := a.getAddr()

	// Initialize the HTTP server
	a.server = &http.Server{
		Addr:    addr,
		Handler: a.router,
	}

	log.Infof("Server started at http://localhost%s%s", addr, a.contextPath)

	// Start the server in a goroutine
	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	return a
}

func (a *apiRouterHandlerImpl[T]) StartServer() {
	if a.port == "" {
		a.port = apiPort()
	}

	if a.contextPath == "" {
		a.contextPath = apiContextPath()
	}

	addr := a.getAddr()
	log.Infof("Server started at http://localhost%s%s", addr, a.contextPath)
	log.Fatal(http.ListenAndServe(addr, a.router))

}

func (a *apiRouterHandlerImpl[T]) getAddr() string {
	addr := ":" + a.port

	if a.port == "80" {
		addr = ""
	} else {
		addr = ":" + a.port
	}
	return addr
}

func (a *apiRouterHandlerImpl[T]) StopServer() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.server == nil {
		return nil // Server is not running
	}

	log.Infof("Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to shut down the server
	if err := a.server.Shutdown(ctx); err != nil {
		log.Errorf("Server shutdown failed: %v", err)
		return err
	}

	log.Println("Server stopped.")
	return nil
}

func (a *apiRouterHandlerImpl[T]) RestartServer() error {
	if err := a.StopServer(); err != nil {
		return err
	}

	// Reinitialize the server
	a.StartServer()
	return nil
}
