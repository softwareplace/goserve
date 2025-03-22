package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// UserHomePathFix resolves a path that starts with '~' to the user's home directory.
// If the home directory cannot be determined, it logs an error and exits the application.
// Parameters:
//   - path: The file path string, potentially starting with '~'.
//
// Returns:
//   - The resolved file path with '~' replaced by the user's home directory, or the original path if no substitution is needed.
func UserHomePathFix(path string) string {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Printf("Error: Unable to resolve user home directory for path: %v", err)
			os.Exit(1)
		}
		path = strings.Replace(path, "~", homeDir, 1)
	}
	return path
}
