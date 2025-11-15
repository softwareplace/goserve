package file

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/softwareplace/goserve/cmd/goserve-generator/config"
	"github.com/softwareplace/goserve/cmd/goserve-generator/template"
	goservestring "github.com/softwareplace/goserve/string"
)

var (
	// JoinPath is a utility variable that points to filepath.Join for constructing file paths in a platform-independent way.
	JoinPath = filepath.Join

	// WriteFile is a variable that points to os.WriteFile, used for writing data to a file, creating it if it doesn't exist.
	WriteFile = os.WriteFile

	// MkdirAll is a variable that points to os.MkdirAll, used for creating directories along with any necessary parents.
	MkdirAll = os.MkdirAll
)

func CreateFile(
	path string,
	content string,
	entries ...goservestring.ReplaceEntry,
) {
	if config.ReplaceCurrent == "false" {
		if _, err := os.Stat(path); err == nil {
			log.Printf("⚠️  File already exists: %s (skipping)", path)
			return
		}
	}

	entries = append(
		entries,
		goservestring.Replacement(template.UsernameKey, config.Username),
		goservestring.Replacement(template.ProjectKey, config.ProjectName),
	)

	for _, entry := range entries {
		if entry.Key != "" || entry.Value != "" {
			content = strings.ReplaceAll(content, entry.Key, entry.Value)
		}
	}

	if err := MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Panicf("❌ Failed to create directories for path %s: %v", path, err)
	}

	if err := WriteFile(path, []byte(content), 0644); err != nil {
		log.Panicf("❌ Failed to create file %s: %v", path, err)
	}
}
