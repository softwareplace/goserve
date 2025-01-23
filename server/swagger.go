package server

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/security/principal"
	httpSwagger "github.com/swaggo/http-swagger"
	"os"
)

func (a *apiRouterHandlerImpl[T]) SetupSwagger(getSwagger func() (swagger *openapi3.T, err error)) ApiRouterHandler[T] {
	swagger, err := getSwagger()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}
	swagger.Servers = nil
	swaggerHandler := httpSwagger.Handler(func(config *httpSwagger.Config) {
		config.URL = ContextPath + "doc"
	})

	a.Router().PathPrefix(ContextPath + "swagger/").Handler(swaggerHandler)

	a.PublicRouter(a.handleSwaggerJSON(swagger), "doc", "GET")
	principal.AddOpenPath("GET::" + ContextPath + "doc")
	principal.AddOpenPath("GET::" + ContextPath + "swagger/.*")
	return a
}

func (a *apiRouterHandlerImpl[T]) handleSwaggerJSON(swagger *openapi3.T) func(ctx *api_context.ApiRequestContext[T]) {
	return func(ctx *api_context.ApiRequestContext[T]) {
		ctx.Response(swagger, 200)
	}
}
