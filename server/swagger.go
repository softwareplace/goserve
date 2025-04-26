package server

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security/router"
	httpSwagger "github.com/swaggo/http-swagger"
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

func (a *baseServer[T]) SwaggerDocProvider(getSwagger func() (swagger *openapi3.T, err error)) Api[T] {
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
		newPath := strings.TrimRight(a.contextPath, "/") + path
		pathLogger(pathItem, newPath)
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
		config.URL = a.contextPath + "doc.json"
		config.Layout = httpSwagger.BaseLayout
	})

	a.Router().PathPrefix(a.contextPath + "swagger/").Handler(swaggerHandler)

	a.PublicRouter(a.handleSwaggerJSON(swagger), "doc.json", "GET")
	router.AddOpenPath("GET::" + a.contextPath + "doc.json")
	router.AddOpenPath("GET::" + a.contextPath + "swagger/.*")
	a.swaggerIsEnabled = true
	return a
}

func (a *baseServer[T]) SwaggerDocHandler(swaggerFile string) Api[T] {
	return a.SwaggerDocProvider(func() (swagger *openapi3.T, err error) {
		return SwaggerDocLoader(swaggerFile)
	})
}

func (a *baseServer[T]) handleSwaggerJSON(swagger *openapi3.T) func(ctx *goservectx.Request[T]) {
	return func(ctx *goservectx.Request[T]) {
		ctx.Response(swagger, 200)
	}
}

func pathLogger(pathItem *openapi3.PathItem, path string) {
	if pathItem.Post != nil {
		log.Printf("POST %s", path)
	}

	if pathItem.Get != nil {
		log.Printf("GET %s", path)
	}

	if pathItem.Put != nil {
		log.Printf("PUT %s", path)
	}

	if pathItem.Delete != nil {
		log.Printf("DELETE %s", path)
	}

	if pathItem.Patch != nil {
		log.Printf("PATCH %s", path)
	}

	if pathItem.Options != nil {
		log.Printf("OPTIONS %s", path)
	}

	if pathItem.Head != nil {
		log.Printf("HEAD %s", path)
	}

	if pathItem.Trace != nil {
		log.Printf("TRACE %s", path)
	}

	if pathItem.Connect != nil {
		log.Printf("CONNECT %s", path)
	}

	if pathItem.Servers != nil {
		log.Printf("SERVERS %s", path)
	}
}
