package main

import (
	"flag"
	"fmt"
	"github.com/softwareplace/goserve/cmd/goserve-generator/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	projectName    = flag.String("n", "", "Project name")
	username       = flag.String("u", "", "GitHub username")
	replaceCurrent = flag.String("r", "false", "Replace current directory/files with generated files")
)

type ReplaceEntry struct {
	Key   string
	Value string
}

func main() {
	flag.Parse()

	if *projectName == "" || *username == "" {
		fmt.Printf("\nOptions:\n")
		fmt.Printf("  -n  (required): Project name (e.g., myproject)\n")
		fmt.Printf("  -u  (required): GitHub username (e.g., myusername)\n")
		fmt.Printf("  -r  (optional): Replace current directory/files with generated files (true/false, default: false)\n")
		fmt.Printf("  -h  (optional): Show this help message\n")
		os.Exit(1)
	}

	dirs := []string{
		"cmd/server",
		"internal/application",
		"internal/adapter/handler",
		"internal/adapter/handler/gen",
		"internal/adapter/repository",
		"internal/adapter/client",
		"internal/core/domain",
		"internal/core/service",
		"internal/core/ports",
		"internal/config",
		"internal/pkg",
		"migrations",
		"config",
		"api",
	}

	root := *projectName
	for _, dir := range dirs {
		path := filepath.Join(root, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatalf("Failed to create %s: %v", path, err)
		}
	}

	createFile(filepath.Join(root, "cmd/server/main.go"), template.GoServeMain)
	createFile(filepath.Join(root, "README.md"), template.Readme)
	createFile(filepath.Join(root, "cmd/server/main_test.go"), template.GoServeMainTest)
	createFile(filepath.Join(root, "internal/adapter/handler/hello.go"), template.Handler)
	createFile(filepath.Join(root, "config/config.yaml"), template.GoServeGenConfig)
	createFile(filepath.Join(root, "api/swagger.yaml"), template.Swagger)
	createFile(filepath.Join(root, "Makefile"), template.Makefile)
	createFile(filepath.Join(root, "go.mod"), template.GoMod)
	createFile(filepath.Join(root, ".gitignore"), template.GitIgnore)
	createFile(filepath.Join(root, "/internal/application/context.go"), template.Context)
	createFile(filepath.Join(root, "/internal/adapter/handler/gen/api.gen.go"), "")

	if err := os.Chdir(root); err != nil {
		log.Fatalf("❌ Failed to change directory to %s: %v", root, err)
	}

	run("make", "test")

	if _, err := os.Stat(root + "/.git"); err != nil {
		run("git", "init")
		run("git", "add", ".")
		run("git", "commit", "-m", "Base project setup")
		run("git", "branch", "-M", "main")
	}

	fmt.Printf("✅ Project %s created successfully!\n", *projectName)

}

func run(command string, args ...string) {
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
	if *replaceCurrent == "false" {
		if _, err := os.Stat(path); err == nil {
			log.Printf("⚠️  File already exists: %s (skipping)", path)
			return
		}
	}

	entries = append(
		entries,
		replacement(template.UsernameKey, *username),
		replacement(template.ProjectKey, *projectName),
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
