package security

import (
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var (
	matcher      = `:[a-zA-Z]+` // Matches dynamic segments like ":param".
	re           = regexp.MustCompile(matcher)
	roles        = make(map[string][]string)
	openPath     []string
	openPathLock sync.RWMutex
)

// AddOpenPath adds a path to the list of open paths.
func AddOpenPath(path string) {
	openPathLock.Lock()
	defer openPathLock.Unlock()
	openPath = append(openPath, path)
}

// AddRoles associates a path with required roles.
func AddRoles(path string, requiredRoles ...string) {
	if len(requiredRoles) > 0 {
		roles[path] = requiredRoles
	}
}

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
