package router

import (
	"regexp"
	"strings"
)

// GetRolesForPath retrieves the roles associated with a request path.
//
// This function takes the API request context and determines the roles required
// for accessing the specified path. The roles are matched based on predefined
// patterns or exact path matches.
//
// Parameters:
//
//	- method: The HTTP method of the request.
//	- path: The path of the request.
//
// Returns:
//
//	- []string: A slice of required roles for the path or nil if no roles are defined.
//	- bool: True if roles are required for the path, false otherwise.
func GetRolesForPath(method, path string) ([]string, bool) {
	resource := method + "::" + path

	for pattern, requiredRoles := range roles {
		regexPattern := convertPathToRegex(pattern)
		regex := regexp.MustCompile(regexPattern)

		if regex.MatchString(resource) || resource == pattern {
			return requiredRoles, true
		}
	}

	return nil, false
}

// IsPublicPath checks if the provided path is registered as a public route.
//
// This function takes the API request context and verifies whether the current
// request path matches any registered public routes. Public routes are those
// that do not require any roles to access.
//
// Parameters:
//	 - method: The HTTP method of the request.
//	 - path: The path of the request.
//
// Returns:
//
//	- bool: True if the path is a public route, false otherwise.
func IsPublicPath(method, path string) bool {
	resource := method + "::" + path
	for _, openPath := range openPaths {
		regexPattern := convertPathToRegex(openPath)
		regex := regexp.MustCompile(regexPattern)
		if regex.MatchString(resource) || resource == openPath {
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
