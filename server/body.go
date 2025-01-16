package server

import (
	"encoding/json"
	"github.com/softwareplace/http-utils/api_context"
	"net/http"
)

type OnSuccess[B any, T api_context.ApiContextData] func(ctx *api_context.ApiRequestContext[T], body B)
type OnError[T api_context.ApiContextData] func(ctx *api_context.ApiRequestContext[T], err error)

func GetRequestBody[B any, T api_context.ApiContextData](
	ctx *api_context.ApiRequestContext[T],
	target B,
	onSuccess OnSuccess[B, T],
	onError OnError[T],
) {
	err := json.NewDecoder(ctx.Request.Body).Decode(&target)
	if err != nil {
		onError(ctx, err)
	} else {
		onSuccess(ctx, target)
	}
}

func FailedToLoadBody[T api_context.ApiContextData](ctx *api_context.ApiRequestContext[T], _ error) {
	ctx.Error("Invalid request data", http.StatusBadRequest)
}
