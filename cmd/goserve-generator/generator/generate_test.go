package generator

import (
	"fmt"
	"github.com/softwareplace/goserve/cmd/goserve-generator/utils"
	testutils "github.com/softwareplace/goserve/internal/utils"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGeneratePathsValidation(t *testing.T) {
	baseProjectPath := utils.JoinPath(testutils.ProjectBasePath(), ".out/test-execution")
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(baseProjectPath)

	t.Run("should create all declared directories and files", func(t *testing.T) {
		require.NotPanics(t, func() {
			Execute(baseProjectPath)
			for _, fileEntry := range generatedFiles {
				require.FileExists(t, utils.JoinPath(baseProjectPath, fileEntry.Path))
			}
		})
	})

	t.Run("should run panic when failed to create a dir", func(t *testing.T) {
		require.NotPanics(t, func() {
			defer func() {
				utils.MkdirAll = os.MkdirAll
			}()

			utils.MkdirAll = func(path string, perm os.FileMode) error {
				return fmt.Errorf("failed to create dir %s, %v", path, perm)
			}
			require.Panics(t, func() {
				createProjectDir(baseProjectPath)
			})
		})
	})
}
