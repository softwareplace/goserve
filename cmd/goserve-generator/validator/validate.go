package validator

import (
	"log"
	"os"
	"os/exec"

	"github.com/softwareplace/goserve/cmd/goserve-generator/cmd"
)

const oapiCodegen = "github.com/deepmap/oapi-codegen/v2/Execute/oapi-codegen@v2.2.0"

var (
	projectExists       = os.Chdir
	cmdMandatoryExecute = cmd.MandatoryExecute
	commandAvailable    = exec.LookPath
)

// ProjectValidate validates and sets up a Go project with API code generation.
// It performs the following tasks:
// - Changes the working directory to the specified project.
// - Ensures the 'oapi-codegen' tool is installed for OpenAPI code generation.
// - Executes code generation, module tidy-up, code formatting, and benchmarking.
func ProjectValidate(projectName string) {
	joinInProject(projectName)
	codeGenValidator()
	codeGenExecute()
}

// joinInProject switches to the specified project directory.
// This ensures that all subsequent file operations are performed
// within the specified project context. If the directory cannot be
// changed, it logs the error and terminates the program.
func joinInProject(projectName string) {
	if err := projectExists(projectName); err != nil {
		log.Panicf("❌ Failed to change directory to %s: %v", projectName, err)
	}
}

// codeGenValidator ensures that the 'oapi-codegen' tool is available for
// OpenAPI code generation. If the tool is not found in the system PATH,
// it installs it automatically using the 'go install' command.
func codeGenValidator() {
	// Check if 'oapi-codegen' is available, if not, install it
	if _, err := commandAvailable("oapi-codegen"); err != nil {
		log.Println("🔍 'oapi-codegen' not found. Installing it...")
		cmdMandatoryExecute("go", "install", oapiCodegen)
		log.Println("✅ 'oapi-codegen' installed successfully.")
	}
}

// codeGenExecute orchestrates the execution of various tasks required to
// set up and validate the project. Tasks include:
// - Generating source code using 'oapi-codegen' based on the OpenAPI spec.
// - Running 'go mod tidy' to synchronize dependencies.
// - Formatting all Go source code files using 'go fmt'.
// - Executing benchmark tests across the project to ensure its stability.
func codeGenExecute() {
	cmdMandatoryExecute("oapi-codegen", "--config", "./config/config.yaml", "./api/swagger.yaml")
	cmdMandatoryExecute("go", "mod", "tidy")
	cmdMandatoryExecute("go", "fmt", "./...")
	cmdMandatoryExecute("go", "test", "-bench=.", "./...")
}
