package git

import (
	"log"
	"os"

	"github.com/softwareplace/goserve/cmd/goserve-generator/cmd"
	"github.com/softwareplace/goserve/cmd/goserve-generator/config"
)

var (
	hasDir         = os.Stat
	runCmd         = cmd.Execute
	gitCommandArgs = [][]string{
		{
			"init", "-q",
		},
		{
			"branch", "-M", "main",
		},
		{
			"add", ".",
		},
		{
			"commit", "-m", "Base project setup",
		},
	}
)

func Setup() {
	if config.GiInit == "true" {
		_, err := hasDir(".git")
		if err != nil {
			log.Printf("Git project initialization: %v", err)
			for _, args := range gitCommandArgs {
				runCmd("git", args...)
			}
		}
	}
}
