package request

import (
	"encoding/json"
	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	"net/http"
	"strings"
)

type OnSuccess[B any, T goservectx.Principal] func(ctx *goservectx.Request[T], body B)
type OnError[T goservectx.Principal] func(ctx *goservectx.Request[T], err error)

func FailedToLoadBody[T goservectx.Principal](ctx *goservectx.Request[T], _ error) {
	ctx.Error("Invalid request data", http.StatusBadRequest)
}

// GetRequestBody parses the JSON request body and executes the appropriate success or error handler.
// ctx is the request context containing headers and the request body.
// target is the variable to decode the request body into.
// onSuccess is invoked if the request body is successfully parsed or if Content-Type is unsupported.
// onError is invoked if JSON decoding fails or any other error occurs.
func GetRequestBody[B any, T goservectx.Principal](
	ctx *goservectx.Request[T],
	target B,
	onSuccess OnSuccess[B, T],
	onError OnError[T],
) {
	goserveerror.Handler(func() {
		// Check if the Content-Type is application/json
		contentType := ctx.Request.Header.Get(goservectx.ContentType)

		if strings.Contains(contentType, goservectx.ApplicationJson) || contentType == "" {
			// Decode the JSON body
			if err := json.NewDecoder(ctx.Request.Body).Decode(&target); err != nil {
				onError(ctx, err)
				return
			}
			onSuccess(ctx, target)
			return
		}

		onSuccess(ctx, target)
	}, func(err error) {
		onError(ctx, err)
	})
}
