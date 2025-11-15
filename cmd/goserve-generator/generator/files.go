package generator

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/softwareplace/goserve/cmd/goserve-generator/config"
	"github.com/softwareplace/goserve/cmd/goserve-generator/template"
	"github.com/softwareplace/goserve/cmd/goserve-generator/version"
	goserverstring "github.com/softwareplace/goserve/string"
)

type generatorFile struct {
	Path    string
	Content string
	Entries []goserverstring.ReplaceEntry
}

// filesGenerator is a collection of predefined files
// with their corresponding paths, content, and optional replacements.
func filesGenerator() []generatorFile {
	return []generatorFile{
		{
			Path:    "cmd/server/main.go",
			Content: template.GoServeMain,
		},
		{
			Path:    "cmd/server/main_test.go",
			Content: template.GoServeMainTest,
		},
		{
			Path:    ".github/workflows/test.yml",
			Content: template.GitHubWorkflow,
		},
		{
			Path:    "internal/adapter/handler/service.go",
			Content: template.HandlerService,
		},
		{
			Path:    "internal/adapter/handler/hello.go",
			Content: template.HandlerImpl,
		},
		{
			Path:    "internal/core/domain/model/model.go",
			Content: template.DomainModel,
		},
		{
			Path:    "internal/application/principal.go",
			Content: template.Context,
		},
		{
			Path:    "internal/application/config/config.go",
			Content: template.AppConfig,
		},
		{
			Path:    "internal/adapter/handler/gen/api.gen.go",
			Content: "",
		},
		{
			Path:    "config/config.yaml",
			Content: getConfigFileContent(),
		},
		{
			Path:    "api/swagger.yaml",
			Content: template.Swagger,
		},
		{
			Path:    "README.md",
			Content: template.Readme,
		},
		{
			Path:    "Makefile",
			Content: template.Makefile,
		},
		{
			Path:    "go.mod",
			Content: template.GoMod,
			Entries: getGoModeReplacementEntries(),
		},
		{
			Path:    ".gitignore",
			Content: template.GitIgnore,
		},
		{
			Path:    "Dockerfile",
			Content: template.Dockerfile,
		},
		{
			Path:    "docker-compose.yaml",
			Content: template.DockerCompose,
		},
	}
}

func getGoModeReplacementEntries() []goserverstring.ReplaceEntry {
	goServeLatest := getGoServeVersion()

	goServerVersion := goserverstring.Replacement(template.GoServeLatestVersionKey, goServeLatest)

	log.Infof("Using custom GoServe version: %s", config.GoServerVersion)
	return []goserverstring.ReplaceEntry{
		goServerVersion,
	}
}

func getGoServeVersion() string {
	var goServeLatest string

	if config.GoServerVersion != "" {
		goServeLatest = config.GoServerVersion
	}

	if goServeLatest == "" {
		goServeLatest = version.GoServeLatest()
	}
	return goServeLatest
}

var (
	readFile = os.ReadFile
)

func getConfigFileContent() string {
	if config.CodeGenConfigFile == "" {
		return template.GoServeGenConfig
	}
	fileContent, err := readFile(config.CodeGenConfigFile)
	if err != nil {
		log.Panicf("❌ Failed to read file %s: %v", config.CodeGenConfigFile, err)

	}
	return string(fileContent)
}
