package context

import (
	"fmt"
	"mime/multipart"
)

// GetSessionId retrieves the unique identifier for the current API session.
// This session ID is used for tracking the lifecycle of requests in a session.
func (ctx *Request[T]) GetSessionId() string {
	return ctx.sessionId
}

// QueryOf retrieves the first value of the specified query parameter from the request URL.
// If the query parameter does not exist or has no values, an empty string is returned.
//
// Parameters:
//   - key: The name of the query parameter to retrieve.
//
// Returns:
//   - The first value of the query parameter or an empty string if it does not exist.
func (ctx *Request[T]) QueryOf(key string) string {
	if len(ctx.QueryValues[key]) > 0 {
		return ctx.QueryValues[key][0]
	}
	return ""
}

// QueriesOf retrieves all values of the specified query parameter from the request URL.
// If the query parameter does not exist, an empty slice is returned.
//
// Parameters:
//   - key: The name of the query parameter to retrieve.
//
// Returns:
//   - A slice of strings containing all values of the query parameter or an empty slice if it does not exist.
func (ctx *Request[T]) QueriesOf(key string) []string {
	return ctx.QueryValues[key]
}

// QueriesOfElse retrieves all values of the specified query parameter from the request URL.
// If the query parameter does not exist, the provided default values are returned.
//
// Parameters:
//   - key: The name of the query parameter to retrieve.
//   - defaultQueries: The default values to return if the query parameter does not exist.
//
// Returns:
//   - A slice of strings containing all values of the query parameter or the default values.
func (ctx *Request[T]) QueriesOfElse(key string, defaultQueries []string) []string {
	if len(ctx.QueryValues[key]) > 0 {
		return ctx.QueryValues[key]
	}
	return defaultQueries
}

// QueryOfOrElse retrieves the first value of the specified query parameter from the request URL.
// If the query parameter does not exist or has no values, the provided default value is returned.
//
// Parameters:
//   - key: The name of the query parameter to retrieve.
//   - defaultQuery: The default value to return if the query parameter does not exist.
//
// Returns:
//   - The first value of the query parameter or the default value.
func (ctx *Request[T]) QueryOfOrElse(key string, defaultQuery string) string {
	if len(ctx.QueryValues[key]) > 0 {
		return ctx.QueryValues[key][0]
	}
	return defaultQuery
}

// HeadersOf retrieves all values of the specified HTTP header from the request.
// If the header does not exist, an empty slice is returned.
//
// Parameters:
//   - key: The name of the HTTP header to retrieve.
//
// Returns:
//   - A slice of strings containing all values of the header or an empty slice if it does not exist.
func (ctx *Request[T]) HeadersOf(key string) []string {
	return ctx.Headers[key]
}

// HeaderOf retrieves the first value of the specified HTTP header from the request.
// If the header does not exist or has no values, an empty string is returned.
//
// Parameters:
//   - key: The name of the HTTP header to retrieve.
//
// Returns:
//   - The first value of the header or an empty string if it does not exist.
func (ctx *Request[T]) HeaderOf(key string) string {
	if len(ctx.Headers[key]) > 0 {
		return ctx.Headers[key][0]
	}
	return ""
}

// HeadersOfOrElse retrieves all values of the specified HTTP header from the request.
// If the header does not exist, the provided default values are returned.
//
// Parameters:
//   - key: The name of the HTTP header to retrieve.
//   - defaultHeaders: The default values to return if the header does not exist.
//
// Returns:
//   - A slice of strings containing all values of the header or the default values.
func (ctx *Request[T]) HeadersOfOrElse(key string, defaultHeaders []string) []string {
	if len(ctx.Headers[key]) > 0 {
		return ctx.Headers[key]
	}
	return defaultHeaders
}

// HeaderOfOrElse retrieves the first value of the specified HTTP header from the request.
// If the header does not exist or has no values, the provided default value is returned.
//
// Parameters:
//   - key: The name of the HTTP header to retrieve.
//   - defaultHeader: The default value to return if the header does not exist.
//
// Returns:
//   - The first value of the header or the default value.
func (ctx *Request[T]) HeaderOfOrElse(key string, defaultHeader string) string {
	if len(ctx.Headers[key]) > 0 {
		return ctx.Headers[key][0]
	}
	return defaultHeader
}

// PathValueOf retrieves the value of the specified path variable from the request URL.
// Path variables are extracted from dynamic segments of the route defined in the router.
//
// Parameters:
//   - key: The name of the path variable to retrieve.
//
// Returns:
//   - The value of the path variable or an empty string if it does not exist.
func (ctx *Request[T]) PathValueOf(key string) string {
	return ctx.PathValues[key]
}

// FormValue retrieves the first value for the given form field name from the parsed form data.
// If the form field does not exist, it returns an empty string.
//
// Parameters:
//   - name: The name of the form field to retrieve.
//
// Returns:
//   - The value of the form field or an empty string if it does not exist.
func (ctx *Request[T]) FormValue(name string) any {
	return ctx.Request.FormValue(name)
}

// FormFile retrieves a file and its header from a multipart form with the given field name.
// The file is immediately closed after being read to avoid resource leaks.
//
// Parameters:
//   - name: The name of the form field containing the file to retrieve.
//
// Returns:
//   - multipart.File: The file object, or nil if an error occurs.
//   - *multipart.FileHeader: The file header, or nil if an error occurs.
//   - error: An error, if one occurs while retrieving the file.
func (ctx *Request[T]) FormFile(name string) (multipart.File, *multipart.FileHeader, error) {
	file, fileHeader, err := ctx.Request.FormFile("resource")
	if err != nil {
		return nil, nil, err
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Failed to close file: %v", err)
		}
	}(file)

	return file, fileHeader, nil
}

// ParseMultipartForm parses the multipart form in the request body, storing up to maxMemory
// bytes of its file parts in memory, with the remainder stored on disk. This is necessary to
// access file uploads sent in a multipart request.
//
// Parameters:
//   - maxMemory: The maximum number of bytes to store in memory.
//
// Returns:
//   - error: An error if parsing fails.
func (ctx *Request[T]) ParseMultipartForm(maxMemory int64) error {
	return ctx.Request.ParseMultipartForm(maxMemory)
}
