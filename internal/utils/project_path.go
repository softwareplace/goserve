package utils

import (
	"fmt"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
)

func TestSecretFilePath() string {
	return fmt.Sprintf("%v/internal/resource/secret/private.key", ProjectBasePath())
}

func ProjectBasePath() string {
	_, b, _, ok := runtime.Caller(0)

	if !ok {
		log.Fatal("Failed to get caller")
	}

	return filepath.Join(filepath.Dir(b), "../..")
}
