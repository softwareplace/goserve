package validator

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/softwareplace/goserve/cmd/goserve-generator/cmd"
	"github.com/softwareplace/goserve/cmd/goserve-generator/config"
	"github.com/softwareplace/goserve/cmd/goserve-generator/generator"
	"github.com/softwareplace/goserve/cmd/goserve-generator/utils"
	testutils "github.com/softwareplace/goserve/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const configFile = `package: gen
generate:
  gorilla-server: true
  models: true
output: ${ROOT_PROJECT}/.out/${PROJECT}/internal/adapter/handler/gen/api.gen.go
output-options:
  user-templates:
    imports.tmpl: ${ROOT_PROJECT}/resource/templates/imports.tmpl
    param-types.tmpl: ${ROOT_PROJECT}/resource/templates/param-types.tmpl
    request-bodies.tmpl: ${ROOT_PROJECT}/resource/templates/request-bodies.tmpl
    typedef.tmpl: ${ROOT_PROJECT}/resource/templates/typedef.tmpl
    gorilla/gorilla-register.tmpl: ${ROOT_PROJECT}/resource/templates/gorilla/gorilla-register.tmpl
    gorilla/gorilla-middleware.tmpl: ${ROOT_PROJECT}/resource/templates/gorilla/gorilla-middleware.tmpl
    gorilla/gorilla-interface.tmpl: ${ROOT_PROJECT}/resource/templates/gorilla/gorilla-interface.tmpl
compatibility:
  apply-gorilla-middleware-first-to-last: true`

func testCleanup(t *testing.T, args ...string) {
	t.Cleanup(func() {
		// Reset global vars after each test
		config.ProjectName = ""
		config.Username = ""
		config.ReplaceCurrent = ""
		config.GiInit = ""
		config.CodeGenConfigFile = ""
		config.GoServerVersion = ""
	})
}

func setArgsForTest(t *testing.T, args ...string) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	err := config.InitFlagsWithSet(fs, args)
	assert.NoError(t, err)
}

func TestValidateProjectValidation(t *testing.T) {
	rootProjectPath := testutils.ProjectBasePath()

	t.Run("should create all declared directories and files", func(t *testing.T) {
		config.ProjectName = "test-execution-validate-01"

		baseProjectPath := utils.JoinPath(rootProjectPath, ".out/", config.ProjectName)
		defer func() {
			_ = os.RemoveAll(baseProjectPath)
			testCleanup(t)
		}()

		setArgsForTest(
			t,
			"-n", config.ProjectName,
			"-u", "test-user",
			"-gi", "false",
			"-r", "true",
			"-gsv", getGitCommitHash(),
		)

		generator.Execute(baseProjectPath)
		configFilePath := utils.JoinPath(baseProjectPath, "config/config.yaml")
		utils.CreateFile(configFilePath, configFile, utils.Replacement("${ROOT_PROJECT}", rootProjectPath))

		ProjectValidate(baseProjectPath)
	})

	t.Run("should exit with panic when project does not exists", func(t *testing.T) {
		config.ProjectName = "test-execution-validate-02"

		baseProjectPath := utils.JoinPath(rootProjectPath, ".out/", config.ProjectName)
		defer func() {
			_ = os.RemoveAll(baseProjectPath)
			testCleanup(t)
		}()

		setArgsForTest(
			t,
			"-n", config.ProjectName,
			"-u", "test-user",
			"-gi", "false",
			"-r", "true",
			"-gsv", getGitCommitHash(),
		)

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
func getGitCommitHash() string {
	// Get current branch name
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOutput, err := branchCmd.Output()
	if err != nil {
		log.Errorf("Failed to get current branch name: %v", err)
		return ""
	}
	branchName := strings.TrimSpace(string(branchOutput))

	// Get local HEAD commit hash
	localCmd := exec.Command("git", "rev-parse", "HEAD")
	localOutput, err := localCmd.Output()
	if err != nil {
		log.Errorf("Failed to get local commit hash: %v", err)
		return ""
	}
	localHash := strings.TrimSpace(string(localOutput))

	// Get remote commit hash (from origin/branch)
	remoteCmd := exec.Command("git", "rev-parse", fmt.Sprintf("origin/%s", branchName))
	remoteOutput, err := remoteCmd.Output()
	if err != nil {
		log.Errorf("Failed to get remote commit hash: %v", err)
		return localHash // Return at least the local hash
	}
	remoteHash := strings.TrimSpace(string(remoteOutput))

	return remoteHash
}
