package git

import (
	"github.com/softwareplace/goserve/cmd/goserve-generator/cmd"
	"github.com/softwareplace/goserve/cmd/goserve-generator/config"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGitSetupValidation(t *testing.T) {
	defer func() {
		hasDir = os.Stat
		runCmd = cmd.Execute
		config.GiInit = "false"
	}()

	t.Run("should run git project initialization when config.GiInit true and has no .git dir", func(t *testing.T) {
		hasDir = func(name string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		}
		config.GiInit = "true"

		var expectCommand [][]string

		runCmd = func(command string, args ...string) {
			fullCommand := append([]string{command}, args...)
			expectCommand = append(expectCommand, fullCommand)
		}

		Setup()

		require.Equalf(
			t,
			len(expectCommand),
			len(gitCommandArgs),
			"Expected %v commands, but got %v",
			len(gitCommandArgs),
			len(expectCommand),
		)

		for i, args := range gitCommandArgs {
			fullCommand := append([]string{"git"}, args...)
			require.Equalf(
				t,
				fullCommand,
				expectCommand[i],
				"Expected command %v, but got %v",
				fullCommand,
				expectCommand[i],
			)
		}
	})

	t.Run("should not run git project initialization when config.GiInit false", func(t *testing.T) {
		config.GiInit = "false"

		var executeCount int

		runCmd = func(command string, args ...string) {
			executeCount++
		}

		Setup()

		require.Equal(t, 0, executeCount)
	})

	t.Run("should not run git project initialization when config.GiInit false but already contains .git dir", func(t *testing.T) {
		config.GiInit = "true"

		hasDir = func(name string) (os.FileInfo, error) {
			return nil, nil
		}

		var executeCount int

		runCmd = func(command string, args ...string) {
			executeCount++
		}

		Setup()

		require.Equal(t, 0, executeCount)
	})
}
