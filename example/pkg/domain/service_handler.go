package domain

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/example/gen"
	"github.com/softwareplace/http-utils/example/pkg/service"
	"github.com/softwareplace/http-utils/server"
	"net/http"
)

type _serverInterfaceDoc struct {
	handler server.ApiRouterHandler[*api_context.DefaultContext]
}
type _serverInterface struct {
	_serverInterfaceDoc
	handler server.ApiRouterHandler[*api_context.DefaultContext]
	service service.Service
}

func ApiRequestHandler(service service.Service, handler server.ApiRouterHandler[*api_context.DefaultContext]) gen.ServerInterface {
	return &_serverInterface{
		handler: handler,
		service: service,
	}
}

// PostLogin operation middleware
func (siw *_serverInterfaceDoc) PostLogin(w http.ResponseWriter, r *http.Request) {

}

// GetTest operation middleware
func (siw *_serverInterface) GetTest(w http.ResponseWriter, r *http.Request) {
	response := siw.service.IsWorking()
	getCtx(w, r).Response(response, 200)
}

// GetTestV2 operation middleware
func (siw *_serverInterface) GetTestV2(w http.ResponseWriter, r *http.Request) {
	response := siw.service.IsWorkingV2()
	getCtx(w, r).Response(response, 200)
}

func getCtx(w http.ResponseWriter, r *http.Request) *api_context.ApiRequestContext[*api_context.DefaultContext] {
	return api_context.Of[*api_context.DefaultContext](w, r, "")
}
