package main

import (
	"fmt"
	"github.com/softwareplace/goserve/cmd/goserve-generator/config"
	"github.com/softwareplace/goserve/cmd/goserve-generator/generator"
	"github.com/softwareplace/goserve/cmd/goserve-generator/git"
	"github.com/softwareplace/goserve/cmd/goserve-generator/validator"
)

func main() {
	config.InitFlags()
	projectName := config.ProjectName
	generator.Execute(projectName)
	validator.ProjectValidate(projectName)
	git.Setup()
	fmt.Printf("✅ Project %s created successfully!\n", config.ProjectName)
}
