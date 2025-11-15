package generator

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/softwareplace/goserve/cmd/goserve-generator/file"
	testutils "github.com/softwareplace/goserve/internal/utils"
)

func TestGeneratePathsValidation(t *testing.T) {
	baseProjectPath := file.JoinPath(testutils.ProjectBasePath(), ".out/generate-test-execution")
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(baseProjectPath)

	t.Run("should create all declared directories and files", func(t *testing.T) {
		require.NotPanics(t, func() {
			Execute(baseProjectPath)
			for _, fileEntry := range filesGenerator() {
				require.FileExists(t, file.JoinPath(baseProjectPath, fileEntry.Path))
			}
		})
	})

	t.Run("should run panic when failed to create a dir", func(t *testing.T) {
		require.NotPanics(t, func() {
			defer func() {
				file.MkdirAll = os.MkdirAll
			}()

			file.MkdirAll = func(path string, perm os.FileMode) error {
				return fmt.Errorf("failed to create dir %s, %v", path, perm)
			}
			require.Panics(t, func() {
				createProjectDir(baseProjectPath)
			})
		})
	})
}
