package validator

import (
	"github.com/softwareplace/goserve/cmd/goserve-generator/cmd"
	"github.com/softwareplace/goserve/cmd/goserve-generator/config"
	"github.com/softwareplace/goserve/cmd/goserve-generator/generator"
	"github.com/softwareplace/goserve/cmd/goserve-generator/utils"
	testutils "github.com/softwareplace/goserve/internal/utils"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"testing"
)

func TestValidateProjectValidation(t *testing.T) {
	baseProjectPath := utils.JoinPath(testutils.ProjectBasePath(), ".out/test-execution")
	defer func(path string) {
		_ = os.RemoveAll(path)
		config.Username = ""
		config.ProjectName = ""
	}(baseProjectPath)

	t.Run("should create all declared directories and files", func(t *testing.T) {
		config.Username = "test-user"
		config.ProjectName = "test-execution"
		require.NotPanics(t, func() {
			generator.Execute(baseProjectPath)
			ProjectValidate(baseProjectPath)
		})
	})

	t.Run("should exit with panic when project does not exists", func(t *testing.T) {
		config.Username = "test-user"
		config.ProjectName = "test-execution"
		projectExists = func(dir string) error {
			return os.ErrNotExist
		}

		require.Panics(t, func() {
			ProjectValidate(baseProjectPath)
		})
	})

	t.Run("should install oapi-codegen command when not available", func(t *testing.T) {

		defer func() {
			commandAvailable = exec.LookPath
			cmdMandatoryExecute = cmd.MandatoryExecute
		}()

		commandAvailable = func(file string) (string, error) {
			return "", os.ErrNotExist
		}

		cmdMandatoryExecuted := false

		cmdMandatoryExecute = func(command string, args ...string) {
			cmdMandatoryExecuted = true
		}

		codeGenValidator()
		require.True(t, cmdMandatoryExecuted)
	})
}
