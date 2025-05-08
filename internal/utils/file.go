package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func ReadFileContent(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Panicf("❌ Failed to read file %s: %v", path, err)
	}
	return string(content)
}
