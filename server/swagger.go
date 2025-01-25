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

func SwaggerDocLoader(swaggerFile string) (swagger *openapi3.T, err error) {

	swagger = &openapi3.T{}
	loader := openapi3.NewLoader()

	file, err := os.Open(swaggerFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open swagger file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close swagger file: %v", err)
		}
	}(file)

	swagger, err = loader.LoadFromFile(swaggerFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse swagger file: %w", err)
	}

	return swagger, nil
}

func (a *apiRouterHandlerImpl[T]) SwaggerDocProvider(getSwagger func() (swagger *openapi3.T, err error)) ApiRouterHandler[T] {
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
		newPath = strings.ReplaceAll(newPath, "{", ":")
		newPath = strings.ReplaceAll(newPath, "}", "")
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
	a.swaggerIsEnabled = true
	return a
}

func (a *apiRouterHandlerImpl[T]) SwaggerDocHandler(swaggerFile string) ApiRouterHandler[T] {
	return a.
		SwaggerDocProvider(func() (swagger *openapi3.T, err error) {
			return SwaggerDocLoader(swaggerFile)
		})
}

func (a *apiRouterHandlerImpl[T]) handleSwaggerJSON(swagger *openapi3.T) func(ctx *api_context.ApiRequestContext[T]) {
	return func(ctx *api_context.ApiRequestContext[T]) {
		ctx.Response(swagger, 200)
	}
}
