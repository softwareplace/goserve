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
	CodeGenConfigFile string
	GoServerVersion   string
	ProjectName       string
	Username          string
	ReplaceCurrent    string
	GiInit            string
	osExit            = os.Exit
	checkVersion      = version.CheckCurrentVersion
	update            = version.Update
)

func argsValidation() {
	args := os.Args

	if len(args) > 1 {
		if args[1] == "version" {
			checkVersion()
			osExit(0)
		}

		if args[1] == "update" {
			update()
			osExit(0)
		}
	}

	flag.Parse()

	if ProjectName == "" || Username == "" {
		flagUsage()
		osExit(1)
	}
}

func InitFlags() {
	flag.StringVar(&ProjectName, "n", "", "Project name")
	flag.StringVar(&Username, "u", "", "GitHub username")
	flag.StringVar(&ReplaceCurrent, "r", "false", "Replace current directory/files with generated files")
	flag.StringVar(&GiInit, "gi", "true", "(optional): Git project initialization")
	flag.StringVar(&CodeGenConfigFile, "cgf", "", "(optional): template of the codegen config file")
	flag.StringVar(&GoServerVersion, "gsv", "", "(optional): use a specific version of goserver")

	flag.Usage = flagUsage

	argsValidation()
}

func flagUsage() {
	_, _ = fmt.Fprintf(os.Stderr, "\nUsage: goserve-generator [options]\n")
	flag.PrintDefaults()
	fmt.Printf("  version\n\tCheck the current version of goserve-generator")
	fmt.Printf("  update\n\tUpdate goserve-generator to the latest version")
}
