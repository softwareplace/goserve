package main

import (
	"flag"
	"fmt"
	"github.com/softwareplace/goserve/cmd/goserve-generator/template"
	"github.com/softwareplace/goserve/cmd/goserve-generator/version"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const oapiCodegen = "github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.2.0"

type ReplaceEntry struct {
	Key   string
	Value string
}

var (
	projectName    string
	username       string
	replaceCurrent string
	giInit         string
	buildVersion   string
)

func init() {
	flag.StringVar(&projectName, "n", "", "Project name")
	flag.StringVar(&username, "u", "", "GitHub username")
	flag.StringVar(&replaceCurrent, "r", "false", "Replace current directory/files with generated files")
	flag.StringVar(&giInit, "gi", "true", "(optional): Git project initialization")

	flag.Usage = func() {
		flagUsage()
	}
}

func flagUsage() {
	_, _ = fmt.Fprintf(os.Stderr, "\nUsage: goserve-generator [options]\n")
	flag.PrintDefaults()
	println("  -version|-v|--version\n\tCheck the current version of goserve-generator")
}

func main() {
	args := os.Args
	if len(args) > 1 && (args[1] == "-version" || args[1] == "-v" || args[1] == "--version") {
		checkCurrentVersion()
		return
	}

	flag.Parse()

	if projectName == "" || username == "" {
		flagUsage()
		os.Exit(1)
	}

	dirs := []string{
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

	root := projectName
	for _, dir := range dirs {
		path := filepath.Join(root, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatalf("Failed to create %s: %v", path, err)
		}
	}

	createFile(filepath.Join(root, "cmd/server/main.go"), template.GoServeMain)
	createFile(filepath.Join(root, ".github/workflows/test.yml"), template.GitHubWorkflow)
	createFile(filepath.Join(root, "README.md"), template.Readme)
	createFile(filepath.Join(root, "cmd/server/main_test.go"), template.GoServeMainTest)
	createFile(filepath.Join(root, "internal/adapter/handler/service.go"), template.HandlerService)
	createFile(filepath.Join(root, "internal/adapter/handler/hello.go"), template.HandlerImpl)
	createFile(filepath.Join(root, "internal/core/domain/model/model.go"), template.DomainModel)
	createFile(filepath.Join(root, "config/config.yaml"), template.GoServeGenConfig)
	createFile(filepath.Join(root, "api/swagger.yaml"), template.Swagger)
	createFile(filepath.Join(root, "Makefile"), template.Makefile)
	createFile(filepath.Join(root, "go.mod"), template.GoMod, replacement(template.GoServeLatestVersionKey, version.GoServeLatest()))
	createFile(filepath.Join(root, ".gitignore"), template.GitIgnore)
	createFile(filepath.Join(root, "internal/application/principal.go"), template.Context)
	createFile(filepath.Join(root, "internal/application/config/config.go"), template.AppConfig)
	createFile(filepath.Join(root, "internal/adapter/handler/gen/api.gen.go"), "")
	createFile(filepath.Join(root, "Dockerfile"), template.Dockerfile)
	createFile(filepath.Join(root, "docker-compose.yaml"), template.DockerCompose)

	if err := os.Chdir(root); err != nil {
		log.Fatalf("❌ Failed to change directory to %s: %v", root, err)
	}

	// Check if 'oapi-codegen' is available, if not, install it
	if _, err := exec.LookPath("oapi-codegen"); err != nil {
		fmt.Println("🔍 'oapi-codegen' not found. Installing it...")
		mandatoryCmd("go", "install", oapiCodegen)
		fmt.Println("✅ 'oapi-codegen' installed successfully.")
	}

	mandatoryCmd("oapi-codegen", "--config", "./config/config.yaml", "./api/swagger.yaml")
	mandatoryCmd("go", "mod", "tidy")
	mandatoryCmd("go", "fmt", "./...")
	mandatoryCmd("go", "test", "-bench=.", "./...")

	if giInit == "true" {
		_, err := os.Stat(".git")
		if err != nil {
			log.Printf("Git project initialization: %v", err)
			cmd("git", "init", "-q")
			cmd("git", "branch", "-M", "main")
			cmd("git", "add", ".")
			cmd("git", "commit", "-m", "Base project setup")
		}
	}

	fmt.Printf("✅ Project %s created successfully!\n", projectName)
}

func checkCurrentVersion() {
	path, err := exec.LookPath("goserve-generator")
	if err != nil {
		log.Fatalf("Could not find goserve-generator: %v", err)
	}

	cmd := exec.Command("go", "version", "-m", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to check version: %v", err)
	}

	// Parse the output to find the version
	version := extractVersion(string(output))
	if version == "" {
		fmt.Println("Could not determine version")
		return
	}

	fmt.Printf("goserve-generator version: %s\n", version)
}

func extractVersion(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "mod\tgithub.com/softwareplace/goserve") || strings.Contains(line, "dep\tgithub.com/softwareplace/goserve") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				return parts[2] // The version is the third field
			}
		}
	}
	return ""
}

func cmd(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

func mandatoryCmd(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ Failed to execute command '%s %v': %v", command, args, err)
	}
}

func replacement(key string, value string) ReplaceEntry {
	return ReplaceEntry{
		Key:   key,
		Value: value,
	}
}

func createFile(path string, content string, entries ...ReplaceEntry) {
	if replaceCurrent == "false" {
		if _, err := os.Stat(path); err == nil {
			log.Printf("⚠️  File already exists: %s (skipping)", path)
			return
		}
	}

	entries = append(
		entries,
		replacement(template.UsernameKey, username),
		replacement(template.ProjectKey, projectName),
	)

	for _, entry := range entries {
		if entry.Key != "" || entry.Value != "" {
			content = strings.ReplaceAll(content, entry.Key, entry.Value)
		}
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		log.Fatalf("❌ Failed to create file %s: %v", path, err)
	}
}
