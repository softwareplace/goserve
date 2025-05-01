package git

import (
	"github.com/softwareplace/goserve/cmd/goserve-generator/cmd"
	"github.com/softwareplace/goserve/cmd/goserve-generator/config"
	"log"
	"os"
)

func Setup() {
	if config.GiInit == "true" {
		_, err := os.Stat(".git")
		if err != nil {
			log.Printf("Git project initialization: %v", err)
			cmd.Execute("git", "init", "-q")
			cmd.Execute("git", "branch", "-M", "main")
			cmd.Execute("git", "add", ".")
			cmd.Execute("git", "commit", "-m", "Base project setup")
		}
	}
}
