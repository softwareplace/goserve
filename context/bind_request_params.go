package context

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/iris-contrib/schema"
	"net/http"
	"net/textproto"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// FieldSource indicates where a field originated from
type FieldSource string

const (
	SourceBody   FieldSource = "body"
	SourceQuery  FieldSource = "query"
	SourcePath   FieldSource = "path"
	SourceHeader FieldSource = "header"
)

// RequestError represents a validation error with contextual information
type RequestError struct {
	Field   string      `json:"field"`   // The original field name from request
	Source  FieldSource `json:"source"`  // Where the field came from
	Message string      `json:"message"` // Human-readable error message
	Code    int         `json:"-"`       // HTTP status code
}

// Error implements the error interface
func (e *RequestError) Error() string {
	return fmt.Sprintf("%s %s", e.Source, e.Message)
}

// BindRequestParams extracts and binds all parameters from the request to the target struct
func (ctx *Request[T]) BindRequestParams(target interface{}) *RequestError {
	r := ctx.Request

	// First bind all parameters
	if err := ctx.bindAllParams(r, target); err != nil {
		return &RequestError{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
	}

	// Then validate required fields from all sources
	if err := ctx.validateAllRequiredFields(target); err != nil {
		return err
	}

	return nil
}

// bindAllParams handles the actual parameter binding
func (ctx *Request[T]) bindAllParams(r *http.Request, target interface{}) error {
	// Convert path vars (map[string]string) to map[string][]string
	pathVars := make(map[string][]string)
	for k, v := range mux.Vars(r) {
		pathVars[k] = []string{v}
	}

	if len(pathVars) > 0 {
		if err := ctx.MapToStruct(pathVars, target, "path"); err != nil {
			return fmt.Errorf("failed to bind path parameters: %w", err)
		}
	}

	// Bind query parameters
	if err := ctx.MapToStruct(r.URL.Query(), target, "query"); err != nil {
		return fmt.Errorf("failed to bind query parameters: %w", err)
	}

	// Convert and bind headers
	headerVars := make(map[string][]string)
	for k, v := range r.Header {
		canonicalName := textproto.CanonicalMIMEHeaderKey(k)
		headerVars[canonicalName] = v
	}

	if err := ctx.MapToStruct(headerVars, target, "header"); err != nil {
		return fmt.Errorf("failed to bind header parameters: %w", err)
	}

	return nil
}

// validateAllRequiredFields checks required fields from all sources
func (ctx *Request[T]) validateAllRequiredFields(target interface{}) *RequestError {
	val := reflect.ValueOf(target)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Determine the source from tags
		var source FieldSource
		switch {
		case field.Tag.Get("path") != "":
			source = SourcePath
		case field.Tag.Get("query") != "":
			source = SourceQuery
		case field.Tag.Get("header") != "":
			source = SourceHeader
		default:
			source = SourceBody
		}

		// Get the original field name from relevant tag
		fieldName := ctx.getFieldNameFromTag(field, source)

		// Check for required tag
		if required, _ := strconv.ParseBool(field.Tag.Get("required")); required {
			if err := ctx.validateField(fieldName, source, fieldVal); err != nil {
				return err
			}
		}
	}

	return nil
}

// getFieldNameFromTag extracts the field name from the appropriate tag
func (ctx *Request[T]) getFieldNameFromTag(field reflect.StructField, source FieldSource) string {
	var tagName string
	switch source {
	case SourcePath:
		tagName = "path"
	case SourceQuery:
		tagName = "query"
	case SourceHeader:
		tagName = "header"
	default:
		tagName = "json"
	}

	if tagValue := field.Tag.Get(tagName); tagValue != "" {
		return strings.Split(tagValue, ",")[0]
	}
	return field.Name
}

// validateField checks if a field meets requirements
func (ctx *Request[T]) validateField(fieldName string, source FieldSource, fieldVal reflect.Value) *RequestError {
	err := &RequestError{
		Field:  fieldName,
		Source: source,
		Code:   http.StatusBadRequest,
	}

	switch fieldVal.Kind() {
	case reflect.String:
		if strings.TrimSpace(fieldVal.String()) == "" {
			err.Message = fmt.Sprintf("parameter [%s] cannot be empty", fieldName)
			return err
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fieldVal.Int() == 0 {
			err.Message = fmt.Sprintf("parameter [%s] is must be non-zero", fieldName)
			return err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if fieldVal.Uint() == 0 {
			err.Message = fmt.Sprintf("parameter [%s] is required and must be positive", fieldName)
			return err
		}
	case reflect.Float32, reflect.Float64:
		if fieldVal.Float() == 0 {
			err.Message = fmt.Sprintf("parameter [%s] is required and must be non-zero", fieldName)
			return err
		}
	case reflect.Ptr, reflect.Interface:
		if fieldVal.IsNil() {
			err.Message = fmt.Sprintf("parameter [%s] is required and must be provided", fieldName)
			return err
		}
	case reflect.Slice, reflect.Array, reflect.Map:
		if fieldVal.Len() == 0 {
			err.Message = fmt.Sprintf("parameter [%s] is required and cannot be empty", fieldName)
			return err
		}
	case reflect.Struct:
		if reflect.DeepEqual(fieldVal.Interface(), reflect.Zero(fieldVal.Type()).Interface()) {
			err.Message = fmt.Sprintf("parameter [%s] is required", fieldName)
			return err
		}
	default:
		return &RequestError{
			Field:   fieldName,
			Source:  source,
			Message: fmt.Sprintf("Unsupported field type '%s' for validation", fieldVal.Kind()),
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}

// MapToStruct binds map values to struct fields using the specified tag
func (ctx *Request[T]) MapToStruct(source map[string][]string, target interface{}, tag string) error {
	decoder := schema.NewDecoder()
	decoder.SetAliasTag(tag)
	decoder.IgnoreUnknownKeys(true)

	values := make(url.Values, len(source))
	for k, v := range source {
		if len(v) > 0 {
			values[k] = v
		}
	}

	return decoder.Decode(target, values)
}
