package server

import (
	"encoding/json"
	"github.com/softwareplace/http-utils/apicontext"
	"net/http"
	"reflect"
	"strings"
)

type OnSuccess[B any, T apicontext.ApiPrincipalContext] func(ctx *apicontext.ApiRequestContext[T], body B)
type OnError[T apicontext.ApiPrincipalContext] func(ctx *apicontext.ApiRequestContext[T], err error)

func FailedToLoadBody[T apicontext.ApiPrincipalContext](ctx *apicontext.ApiRequestContext[T], _ error) {
	ctx.Error("Invalid request data", http.StatusBadRequest)
}

func GetRequestBody[B any, T apicontext.ApiPrincipalContext](
	ctx *apicontext.ApiRequestContext[T],
	target B,
	onSuccess OnSuccess[B, T],
	onError OnError[T],
) {
	// Check if the Content-Type is application/json
	contentType := ctx.Request.Header.Get("Content-Type")

	if !strings.Contains(contentType, "application/json") {
		onSuccess(ctx, target)
		return
	}

	// Decode the JSON body
	err := json.NewDecoder(ctx.Request.Body).Decode(&target)
	if err != nil {
		onError(ctx, err)
	} else {
		onSuccess(ctx, target)
	}
}

func PopulateFieldsFromRequest[B any, T apicontext.ApiPrincipalContext](
	ctx *apicontext.ApiRequestContext[T],
	target *B, // Pass a pointer to target
) {
	// Use reflection to iterate over the target's fields
	v := reflect.ValueOf(target).Elem() // Dereference the pointer to get the struct value
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Extract field information
		fieldValue := v.Field(i)
		if !fieldValue.CanSet() {
			continue // Skip unexported fields
		}

		// Use the json tag to map the field
		fieldTag := field.Tag.Get("json")
		var tValue any

		// Match value from Path, Query, or Headers
		if ctx.PathValues[fieldTag] != "" {
			tValue = ctx.PathValues[fieldTag]
		} else if ctx.QueryValues[fieldTag] != nil && len(ctx.QueryValues[fieldTag]) > 0 {
			tValue = ctx.QueryValues[fieldTag]
		} else if ctx.Headers[fieldTag] != nil && len(ctx.Headers[fieldTag]) > 0 {
			tValue = ctx.Headers[fieldTag][0]
		}

		// If a value was found, set it to the field
		if tValue != nil {
			// Convert the type of tValue to match the field type
			value := reflect.ValueOf(tValue)
			converted := value.Convert(fieldValue.Type())
			fieldValue.Set(converted)
		}
	}
}
