package shared

import (
	"github.com/softwareplace/http-utils/api_context"
	"regexp"
	"strings"
)

// GetRolesForPath retrieves the roles associated with a request path.
func GetRolesForPath[T api_context.ApiContextData](ctx api_context.ApiRequestContext[T]) ([]string, bool) {
	path := ctx.Request.Method + "::" + ctx.Request.URL.Path

	for pattern, requiredRoles := range roles {
		regexPattern := convertPathToRegex(pattern)
		regex := regexp.MustCompile(regexPattern)

		if regex.MatchString(path) || path == pattern {
			return requiredRoles, true
		}
	}

	return nil, false
}

// IsPublicPath  checks if the provided path is registered as a public route (does not require any roles).
func IsPublicPath[T api_context.ApiContextData](ctx api_context.ApiRequestContext[T]) bool {
	path := ctx.Request.Method + "::" + ctx.Request.URL.Path
	for _, openPath := range openPaths {
		regexPattern := convertPathToRegex(openPath)
		regex := regexp.MustCompile(regexPattern)
		if regex.MatchString(path) || path == openPath {
			return true
		}
	}
	return false
}

// convertPathToRegex converts a path with dynamic segments (e.g., ":param") into a regex pattern.
func convertPathToRegex(path string) string {
	// Escape slashes and replace dynamic segments with regex groups.
	escapedPath := strings.ReplaceAll(path, "/", `\/`)
	return "^" + re.ReplaceAllString(escapedPath, `[^\/]+`) + "$"
}
