package server

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/security/principal"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"os"
	"strings"
)

func (a *apiRouterHandlerImpl[T]) SetupSwagger(getSwagger func() (swagger *openapi3.T, err error)) ApiRouterHandler[T] {
	swagger, err := getSwagger()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}
	swagger.Servers = nil

	// Dereference swagger.Paths to iterate over the map
	// Copy swagger.Paths to a new variable
	paths := make(map[string]*openapi3.PathItem)
	var oldsPaths []string
	for path, pathItem := range swagger.Paths.Map() {
		newPath := strings.TrimRight(ContextPath, "/") + path
		log.Printf("path: %s", newPath)
		oldsPaths = append(oldsPaths, path)
		paths[newPath] = pathItem
	}

	for _, e := range oldsPaths {
		swagger.Paths.Delete(e)
	}

	for path, pathItem := range paths {
		swagger.Paths.Set(path, pathItem)
	}

	swaggerHandler := httpSwagger.Handler(func(config *httpSwagger.Config) {
		config.URL = ContextPath + "doc.json"
		config.Layout = httpSwagger.BaseLayout
	})

	a.Router().PathPrefix(ContextPath + "swagger/").Handler(swaggerHandler)

	a.PublicRouter(a.handleSwaggerJSON(swagger), "doc.json", "GET")
	principal.AddOpenPath("GET::" + ContextPath + "doc.json")
	principal.AddOpenPath("GET::" + ContextPath + "swagger/.*")
	return a
}

func (a *apiRouterHandlerImpl[T]) handleSwaggerJSON(swagger *openapi3.T) func(ctx *api_context.ApiRequestContext[T]) {
	return func(ctx *api_context.ApiRequestContext[T]) {
		ctx.Response(swagger, 200)
	}
}
