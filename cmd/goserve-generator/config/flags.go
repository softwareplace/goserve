package config

import (
	"flag"
	"fmt"
	"github.com/softwareplace/goserve/cmd/goserve-generator/version"
	"log"
	"os"
)

// disable log do add date on log
func init() {
	log.SetFlags(0)

}

var (
	ProjectName    string
	Username       string
	ReplaceCurrent string
	GiInit         string
)

func argsValidation() {
	args := os.Args

	if len(args) > 1 {
		if args[1] == "version" {
			version.CheckCurrentVersion()
			os.Exit(0)
		}

		if args[1] == "update" {
			version.Update()
			os.Exit(0)
		}
	}

	flag.Parse()

	if ProjectName == "" || Username == "" {
		flagUsage()
		os.Exit(1)
	}
}

func InitFlags() {
	flag.StringVar(&ProjectName, "n", "", "Project name")
	flag.StringVar(&Username, "u", "", "GitHub username")
	flag.StringVar(&ReplaceCurrent, "r", "false", "Replace current directory/files with generated files")
	flag.StringVar(&GiInit, "gi", "true", "(optional): Git project initialization")

	flag.Usage = func() {
		flagUsage()
	}

	argsValidation()
}

func flagUsage() {
	_, _ = fmt.Fprintf(os.Stderr, "\nUsage: goserve-generator [options]\n")
	flag.PrintDefaults()
	fmt.Printf("  version\n\tCheck the current version of goserve-generator")
	fmt.Printf("  update\n\tUpdate goserve-generator to the latest version")
}
