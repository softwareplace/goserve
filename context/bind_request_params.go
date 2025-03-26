package context

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/softwareplace/goserve/reflect"
	"net/http"
)

// FieldSource indicates where a field originated from
type FieldSource string

// RequestError represents a validation error with contextual information
type RequestError struct {
	Field   string      `json:"field"`      // The original field name from request
	Source  FieldSource `json:"source"`     // Where the field came from
	Message string      `json:"message"`    // Human-readable error message
	Code    int         `json:"statusCode"` // HTTP status code
}

// Error implements the error interface
func (e *RequestError) Error() string {
	return fmt.Sprintf("%s %s", e.Source, e.Message)
}

// BindRequestParams extracts and binds all parameters from the request to the target struct
func (ctx *Request[T]) BindRequestParams(target interface{}) *RequestError {
	r := ctx.Request

	_ = reflect.ParamsExtract(target,
		reflect.ParamsExtractorSource{
			Tree: r.URL.Query(),
		}, reflect.ParamsExtractorSource{
			Tree: r.Header,
		}, reflect.ParamsExtractorSource{
			Source: mux.Vars(r),
		},
	)

	err := ctx.StructValidation(target)

	if err != nil {
		return &RequestError{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
	}
	return nil
}
