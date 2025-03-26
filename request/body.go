package request

import (
	"encoding/json"
	apicontext "github.com/softwareplace/goserve/context"
	"net/http"
	"strings"
)

type OnSuccess[B any, T apicontext.Principal] func(ctx *apicontext.Request[T], body B)
type OnError[T apicontext.Principal] func(ctx *apicontext.Request[T], err error)

func FailedToLoadBody[T apicontext.Principal](ctx *apicontext.Request[T], _ error) {
	ctx.Error("Invalid request data", http.StatusBadRequest)
}

// GetRequestBody parses the JSON request body and executes the appropriate success or error handler.
// ctx is the request context containing headers and the request body.
// target is the variable to decode the request body into.
// onSuccess is invoked if the request body is successfully parsed or if Content-Type is unsupported.
// onError is invoked if JSON decoding fails or any other error occurs.
func GetRequestBody[B any, T apicontext.Principal](
	ctx *apicontext.Request[T],
	target B,
	onSuccess OnSuccess[B, T],
	onError OnError[T],
) {
	// Check if the Content-Type is application/json
	contentType := ctx.Request.Header.Get(apicontext.ContentType)

	if strings.Contains(contentType, apicontext.ApplicationJson) || contentType == "" {
		// Decode the JSON body
		if err := json.NewDecoder(ctx.Request.Body).Decode(&target); err != nil {
			onError(ctx, err)
			return
		}
		onSuccess(ctx, target)
		return
	}

	onSuccess(ctx, target)
}
