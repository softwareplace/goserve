package generator

import (
	"log"

	"github.com/softwareplace/goserve/cmd/goserve-generator/utils"
)

var (
	dirs = []string{
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
)

// Execute is the main entry point for generating a complete project structure,
// including directories, filesGenerator, and configuration templates.
//
// Responsibilities of Execute:
// 1. Create the project directory structure using predefined templates.
// 2. Generate essential project filesGenerator (e.g., main.go, GitHub workflows).
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
// - projectName: The name of the project for which filesGenerator will be generated.
func createBaseProjectDirAndFiles(projectName string) {
	for _, fileEntry := range filesGenerator() {
		utils.CreateFile(
			utils.JoinPath(projectName, fileEntry.Path),
			fileEntry.Content,
			fileEntry.Entries...,
		)
	}
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
		path := utils.JoinPath(projectName, dir)
		if err := utils.MkdirAll(path, 0755); err != nil {
			log.Panicf("❌ Failed to create directory %s: %v", path, err)
		}
	}
}
