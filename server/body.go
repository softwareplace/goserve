package server

import (
	"encoding/json"
	"net/http"
)

type OnSuccess[T any] func(ctx *ApiRequestContext, body T)
type OnError func(ctx *ApiRequestContext, err error)

func GetRequestBody[T any](ctx *ApiRequestContext, target T, onSuccess OnSuccess[T], onError OnError) {
	err := json.NewDecoder(ctx.Request.Body).Decode(&target)
	if err != nil {
		onError(ctx, err)
	} else {
		onSuccess(ctx, target)
	}
}

func FailedToLoadBody(ctx *ApiRequestContext, _ error) {
	ctx.Error("Invalid request data", http.StatusBadRequest)
}
