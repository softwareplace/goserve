package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
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
