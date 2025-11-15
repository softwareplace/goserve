package generator

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/softwareplace/goserve/cmd/goserve-generator/config"
	"github.com/softwareplace/goserve/cmd/goserve-generator/template"
	"github.com/softwareplace/goserve/cmd/goserve-generator/utils"
	"github.com/softwareplace/goserve/cmd/goserve-generator/version"
	testutils "github.com/softwareplace/goserve/internal/utils"
)

const configFile = `package: gen
generate:
  gorilla-server: true
  models: true
output: ${ROOT_PROJECT}/.out/file-config-test-execution/internal/adapter/handler/gen/api.gen.go
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

func TestGetTargetVersion(t *testing.T) {
	t.Run("should return default goserve version when config.GoServerVersion was not set", func(t *testing.T) {
		defaultGoServeLatest := version.GoServeLatest
		defer func() {
			version.GoServeLatest = defaultGoServeLatest
		}()

		version.GoServeLatest = func() string {
			return "1.0.0-test-only"
		}

		require.Equal(t, "1.0.0-test-only", getGoServeVersion())
	})

	t.Run("should return provide goserve version when config.GoServerVersion was set", func(t *testing.T) {
		defer func() {
			config.GoServerVersion = ""
		}()
		config.GoServerVersion = "1.0.0-test-only"
		require.Equal(t, "1.0.0-test-only", getGoServeVersion())
	})
}

func TestGetFileContent(t *testing.T) {
	t.Run("should return default goserve gen config file wen config.CodeGenConfigFile was not set", func(t *testing.T) {
		require.Equal(t, template.GoServeGenConfig, getConfigFileContent())
	})

	t.Run("should return provided goserve gen config file wen config.CodeGenConfigFile was set", func(t *testing.T) {
		config.ProjectName = "file-config-test-execution"
		baseProjectPath := utils.JoinPath(testutils.ProjectBasePath(), ".out/"+config.ProjectName)
		defer func(path string) {
			_ = os.RemoveAll(path)
			config.CodeGenConfigFile = ""
			config.ProjectName = ""
		}(baseProjectPath)

		configPath := utils.JoinPath(baseProjectPath, "config.yaml")

		utils.CreateFile(configPath, configFile)

		config.CodeGenConfigFile = configPath

		require.Equal(t, configFile, getConfigFileContent())
	})

	t.Run("should exit with panic when provided config file does not exists", func(t *testing.T) {
		defer func() {
			config.CodeGenConfigFile = ""
		}()

		config.CodeGenConfigFile = "/invalid/path/config.yaml"

		require.Panics(t, func() {
			getConfigFileContent()
		})
	})
}
