package api_context

// GetSessionId retrieves the unique identifier for the current API session.
// This session ID is used for tracking the lifecycle of requests in a session.
func (ctx *ApiRequestContext[T]) GetSessionId() string {
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
func (ctx *ApiRequestContext[T]) QueryOf(key string) string {
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
func (ctx *ApiRequestContext[T]) QueriesOf(key string) []string {
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
func (ctx *ApiRequestContext[T]) QueriesOfElse(key string, defaultQueries []string) []string {
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
func (ctx *ApiRequestContext[T]) QueryOfOrElse(key string, defaultQuery string) string {
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
func (ctx *ApiRequestContext[T]) HeadersOf(key string) []string {
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
func (ctx *ApiRequestContext[T]) HeaderOf(key string) string {
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
func (ctx *ApiRequestContext[T]) HeadersOfOrElse(key string, defaultHeaders []string) []string {
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
func (ctx *ApiRequestContext[T]) HeaderOfOrElse(key string, defaultHeader string) string {
	if len(ctx.Headers[key]) > 0 {
		return ctx.Headers[key][0]
	}
	return defaultHeader
}

// PathValuesOf retrieves the value of the specified path variable from the request URL.
// Path variables are extracted from dynamic segments of the route defined in the router.
//
// Parameters:
//   - key: The name of the path variable to retrieve.
//
// Returns:
//   - The value of the path variable or an empty string if it does not exist.
func (ctx *ApiRequestContext[T]) PathValuesOf(key string) string {
	return ctx.PathValues[key]
}
