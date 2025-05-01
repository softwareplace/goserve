package generator

import (
	"github.com/softwareplace/goserve/cmd/goserve-generator/template"
	"github.com/softwareplace/goserve/cmd/goserve-generator/utils"
	"github.com/softwareplace/goserve/cmd/goserve-generator/version"
	"log"
	"os"
	"path/filepath"
)

var dirs = []string{
	"cmd/server",
	".github/workflows",
	"internal/application",
	"internal/application/config",
	"internal/adapter/handler",
	"internal/adapter/handler/gen",
	"internal/adapter/repository",
	"internal/adapter/client",
	"internal/core/domain",
	"internal/core/domain/model",
	"internal/core/service",
	"internal/core/ports",
	"internal/pkg",
	"migrations",
	"config",
	"api",
}

// Execute is the main entry point for generating a complete project structure,
// including directories, files, and configuration templates.
//
// Responsibilities of Execute:
// 1. Create the project directory structure using predefined templates.
// 2. Generate essential project files (e.g., main.go, GitHub workflows).
// 3. Switch to the project directory for subsequent operations.
// 4. Validate and ensure installation of external dependencies (e.g., 'oapi-codegen').
// 5. Perform auxiliary operations such as module tidy, formatting, and testing.
//
// Parameters:
// - projectName: The name of the project to be created.
func Execute(projectName string) {
	// Create the project directory structure
	createProjectDir(projectName)
	createBaseProjectDirAndFiles(projectName)
}

// createBaseProjectDirAndFiles is responsible for creating and populating essential
// project files based on predefined templates.
//
// Responsibilities:
// - Generate main application entry files (e.g., main.go, main_test.go).
// - Set up CI/CD workflows like GitHub Actions.
// - Populate internal folders with core logic (handlers, models, services, etc.).
// - Create configuration files (e.g., YAML, Swagger API spec).
// - Generate project metadata files (e.g., README.md, go.mod, Makefile).
//
// Parameters:
// - projectName: The name of the project for which files will be generated.
func createBaseProjectDirAndFiles(projectName string) {
	// Create main application files
	utils.CreateFile(filepath.Join(projectName, "cmd/server/main.go"), template.GoServeMain)
	utils.CreateFile(filepath.Join(projectName, "cmd/server/main_test.go"), template.GoServeMainTest)

	// Setup GitHub workflows
	utils.CreateFile(filepath.Join(projectName, ".github/workflows/test.yml"), template.GitHubWorkflow)

	// Generate internal structures and logic template files
	utils.CreateFile(filepath.Join(projectName, "internal/adapter/handler/service.go"), template.HandlerService)
	utils.CreateFile(filepath.Join(projectName, "internal/adapter/handler/hello.go"), template.HandlerImpl)
	utils.CreateFile(filepath.Join(projectName, "internal/core/domain/model/model.go"), template.DomainModel)
	utils.CreateFile(filepath.Join(projectName, "internal/application/principal.go"), template.Context)
	utils.CreateFile(filepath.Join(projectName, "internal/application/config/config.go"), template.AppConfig)
	utils.CreateFile(filepath.Join(projectName, "internal/adapter/handler/gen/api.gen.go"), "")

	// Create configuration, API, and documentation files
	utils.CreateFile(filepath.Join(projectName, "config/config.yaml"), template.GoServeGenConfig)
	utils.CreateFile(filepath.Join(projectName, "api/swagger.yaml"), template.Swagger)
	utils.CreateFile(filepath.Join(projectName, "README.md"), template.Readme)
	utils.CreateFile(filepath.Join(projectName, "Makefile"), template.Makefile)

	// Create project metadata and deployment files
	utils.CreateFile(filepath.Join(projectName, "go.mod"), template.GoMod, utils.Replacement(
		template.GoServeLatestVersionKey,
		version.GoServeLatest(),
	))
	utils.CreateFile(filepath.Join(projectName, ".gitignore"), template.GitIgnore)
	utils.CreateFile(filepath.Join(projectName, "Dockerfile"), template.Dockerfile)
	utils.CreateFile(filepath.Join(projectName, "docker-compose.yaml"), template.DockerCompose)
}

// createProjectDir creates the directory structure for the project. It uses a predefined
// list of directory paths (stored in the global 'dirs' variable) and ensures each directory
// is created under the provided project name.
//
// If the creation of any directory fails, the function logs the error and terminates execution.
//
// Parameters:
// - projectName: The name of the project whose directory structure will be created.
func createProjectDir(projectName string) {
	for _, dir := range dirs {
		path := filepath.Join(projectName, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatalf("❌ Failed to create directory %s: %v", path, err)
		}
	}
}
