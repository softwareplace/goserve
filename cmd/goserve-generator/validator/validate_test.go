package validator

import (
	log "github.com/sirupsen/logrus"
	"github.com/softwareplace/goserve/cmd/goserve-generator/cmd"
	"github.com/softwareplace/goserve/cmd/goserve-generator/config"
	"github.com/softwareplace/goserve/cmd/goserve-generator/generator"
	"github.com/softwareplace/goserve/cmd/goserve-generator/utils"
	testutils "github.com/softwareplace/goserve/internal/utils"
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

func TestValidateProjectValidation(t *testing.T) {
	config.GiInit = "false"
	rootProjectPath := testutils.ProjectBasePath()

	t.Run("should create all declared directories and files", func(t *testing.T) {
		config.ProjectName = "test-execution-validate-01"

		baseProjectPath := utils.JoinPath(rootProjectPath, ".out/", config.ProjectName)
		defer func() {
			_ = os.RemoveAll(baseProjectPath)
			config.Username = ""
			config.ProjectName = ""
			config.GoServerVersion = ""
		}()

		config.Username = "test-user"
		config.GoServerVersion = getGitCommitHash()

		generator.Execute(baseProjectPath)
		utils.CreateFile(utils.JoinPath(baseProjectPath, "config/config.yaml"), configFile, utils.Replacement("${ROOT_PROJECT}", rootProjectPath))

		log.Infof("go.mod \n%s", testutils.ReadFileContent(utils.JoinPath(baseProjectPath, "go.mod")))

		ProjectValidate(baseProjectPath)
	})

	t.Run("should exit with panic when project does not exists", func(t *testing.T) {
		config.ProjectName = "test-execution-validate-02"

		baseProjectPath := utils.JoinPath(rootProjectPath, ".out/", config.ProjectName)
		defer func() {
			_ = os.RemoveAll(baseProjectPath)
			config.Username = ""
			config.ProjectName = ""
			config.GoServerVersion = ""
		}()

		config.Username = "test-user"
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
	command := exec.Command("git", "rev-parse", "HEAD")
	output, err := command.Output()
	if err != nil {
		log.Errorf("Failed to get git commit hash: %v", err)
	}
	return strings.TrimSpace(string(output))
}
