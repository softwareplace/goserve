package utils

import (
	"github.com/softwareplace/goserve/cmd/goserve-generator/config"
	"github.com/softwareplace/goserve/cmd/goserve-generator/template"
	"log"
	"os"
	"strings"
)

type ReplaceEntry struct {
	Key   string
	Value string
}

func CreateFile(
	path string,
	content string,
	entries ...ReplaceEntry,
) {
	if config.ReplaceCurrent == "false" {
		if _, err := os.Stat(path); err == nil {
			log.Printf("⚠️  File already exists: %s (skipping)", path)
			return
		}
	}

	entries = append(
		entries,
		Replacement(template.UsernameKey, config.Username),
		Replacement(template.ProjectKey, config.ProjectName),
	)

	for _, entry := range entries {
		if entry.Key != "" || entry.Value != "" {
			content = strings.ReplaceAll(content, entry.Key, entry.Value)
		}
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		log.Fatalf("❌ Failed to create file %s: %v", path, err)
	}
}
