package generator

import (
	"github.com/softwareplace/goserve/cmd/goserve-generator/template"
	"github.com/softwareplace/goserve/cmd/goserve-generator/utils"
	"github.com/softwareplace/goserve/cmd/goserve-generator/version"
)

// generatedFiles is a collection of predefined files
// with their corresponding paths, content, and optional replacements.
var generatedFiles = []struct {
	Path    string
	Content string
	Entries []utils.ReplaceEntry
}{
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
		Content: template.GoServeGenConfig,
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
		Entries: []utils.ReplaceEntry{
			utils.Replacement(template.GoServeLatestVersionKey, version.GoServeLatest()),
		},
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
