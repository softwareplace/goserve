package server

import (
	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/http"
	"github.com/softwareplace/goserve/security/model"
)

// ApiKeyGenerator handles the generation of API keys by processing the request body
// and delegating to the secret service handler. It ensures proper error handling
// for cases where the request body cannot be loaded.
func (a *baseServer[T]) ApiKeyGenerator(ctx *goservectx.Request[T]) {
	http.GetRequestBody(ctx, model.ApiKeyEntryData{}, a.secretService.Handler, http.FailedToLoadBody[T])
}
