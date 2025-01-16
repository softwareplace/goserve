package shared

import (
	"net/http"
	"regexp"
	"strings"
)

// GetRolesForPath retrieves the roles associated with a request path.
func GetRolesForPath(r *http.Request) ([]string, bool) {
	path := r.Method + "::" + r.URL.Path

	for pattern, requiredRoles := range roles {
		regexPattern := convertPathToRegex(pattern)
		regex := regexp.MustCompile(regexPattern)

		if regex.MatchString(path) || path == pattern {
			return requiredRoles, true
		}
	}

	return nil, false
}

// convertPathToRegex converts a path with dynamic segments (e.g., ":param") into a regex pattern.
func convertPathToRegex(path string) string {
	// Escape slashes and replace dynamic segments with regex groups.
	escapedPath := strings.ReplaceAll(path, "/", `\/`)
	return "^" + re.ReplaceAllString(escapedPath, `[^\/]+`) + "$"
}
