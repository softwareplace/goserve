package server

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/error_handler"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/security/principal"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type apiRouterHandlerImpl[T api_context.ApiPrincipalContext] struct {
	router                              *mux.Router
	principalService                    principal.PService[T]
	errorHandler                        error_handler.ApiErrorHandler[T]
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

func (a *apiRouterHandlerImpl[T]) NotFoundHandler() ApiRouterHandler[T] {
	a.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("404 page not found: %s", r.URL.Path)

		if r.Method == "GET" {
			swaggerPath := strings.TrimSuffix(a.contextPath, "/") + "/swagger"

			isSwaggerPath := strings.TrimSuffix(r.URL.Path, "/") == swaggerPath

			if a.swaggerIsEnabled && (r.URL.Path == a.contextPath || r.URL.Path == a.contextPath[:len(a.contextPath)-1] || isSwaggerPath) {
				a.goToSwaggerUi(w, r)
				return
			}
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
	//if a.port == "" {
	//	a.port = apiPort()
	//}
	//
	//if a.contextPath == "" {
	//	a.contextPath = apiContextPath()
	//}
	//
	//log.Printf("Server started at http://localhost:%s%s", a.port, a.contextPath)
	//log.Fatal(http.ListenAndServe(":"+a.port, a.router))
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.port == "" {
		a.port = apiPort()
	}

	if a.contextPath == "" {
		a.contextPath = apiContextPath()
	}

	// Initialize the HTTP server
	a.server = &http.Server{
		Addr:    ":" + a.port,
		Handler: a.router,
	}

	log.Printf("Server started at http://localhost:%s%s", a.port, a.contextPath)

	// Start the server in a goroutine
	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()
}

func (a *apiRouterHandlerImpl[T]) StopServer() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.server == nil {
		return nil // Server is not running
	}

	log.Println("Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to shut down the server
	if err := a.server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
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
