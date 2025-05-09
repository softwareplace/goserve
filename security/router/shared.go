package router

import (
	"regexp"
	"sync"
)

var (
	matcher      = `:[a-zA-Z]+` // Matches dynamic segments like ":param".
	re           = regexp.MustCompile(matcher)
	roles        = make(map[string][]string)
	openPaths    []string
	openPathLock sync.RWMutex
)

// AddOpenPath adds a path to the list of open paths.
func AddOpenPath(path string) {

	path = regexp.MustCompile(`/+`).ReplaceAllString(path, "/")
	openPathLock.Lock()
	defer openPathLock.Unlock()
	for _, existingPath := range openPaths {
		if existingPath == path {
			return
		}
	}
	openPaths = append(openPaths, path)
}

// AddRoles associates a path with required roles.
func AddRoles(path string, requiredRoles ...string) {
	if len(requiredRoles) > 0 {
		roles[path] = requiredRoles
	}
}
