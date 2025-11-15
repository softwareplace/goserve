package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/softwareplace/goserve/cmd/goserve-generator/config"
	testutils "github.com/softwareplace/goserve/internal/utils"
)

func TestGeneratePathsValidation(t *testing.T) {
	baseProjectPath := JoinPath(testutils.ProjectBasePath(), ".out/test-execution")
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(baseProjectPath)

	t.Run("should create all declared directories and files", func(t *testing.T) {
		CreateFile(JoinPath(baseProjectPath, "config_ReplaceCurrent_test01.txt"), "test content")
		require.FileExists(t, JoinPath(baseProjectPath, "config_ReplaceCurrent_test01.txt"))
	})

	t.Run("should should not replace file content when already exists and config.ReplaceCurrent false", func(t *testing.T) {
		config.ReplaceCurrent = "false"

		defer func() {
			config.ReplaceCurrent = "true"
		}()

		CreateFile(JoinPath(baseProjectPath, "config_ReplaceCurrent_test02.txt"), "test content 01")
		require.FileExists(t, JoinPath(baseProjectPath, "config_ReplaceCurrent_test02.txt"))
		CreateFile(JoinPath(baseProjectPath, "config_ReplaceCurrent_test02.txt"), "test content 02")

		require.Equal(t, "test content 01", testutils.ReadFileContent(JoinPath(baseProjectPath, "config_ReplaceCurrent_test02.txt")))
	})

	t.Run("should should replace file content when already exists and config.ReplaceCurrent true", func(t *testing.T) {
		config.ReplaceCurrent = "true"
		CreateFile(JoinPath(baseProjectPath, "config_ReplaceCurrent_test03.txt"), "test content 01")
		require.FileExists(t, JoinPath(baseProjectPath, "config_ReplaceCurrent_test03.txt"))
		CreateFile(JoinPath(baseProjectPath, "config_ReplaceCurrent_test03.txt"), "test content 02")

		require.Equal(t, "test content 02", testutils.ReadFileContent(JoinPath(baseProjectPath, "config_ReplaceCurrent_test03.txt")))
	})

	t.Run("should run panic when failed to create a dir", func(t *testing.T) {
		require.NotPanics(t, func() {
			defer func() {
				MkdirAll = os.MkdirAll
			}()

			MkdirAll = func(path string, perm os.FileMode) error {
				return fmt.Errorf("failed to create dir %s, %v", path, perm)
			}
			require.Panics(t, func() {
				CreateFile("config_ReplaceCurrent_test04.txt", "test content")
			})
		})
	})

	t.Run("should run panic when failed to create a file", func(t *testing.T) {
		require.NotPanics(t, func() {
			defer func() {
				WriteFile = os.WriteFile
			}()

			WriteFile = func(path string, data []byte, perm os.FileMode) error {
				return fmt.Errorf("failed to create file %s, %v", path, perm)
			}

			require.Panics(t, func() {
				CreateFile("config_ReplaceCurrent_test04.txt", "test content")
			})
		})
	})
}
