package version

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	executableName        = "goserve-generator"
	checkVersion          = CheckCurrentVersion
	getLatestVersion      = GoServeLatest
	extractCurrentVersion = extractVersion
	gitTargetInstaller    = "github.com/softwareplace/goserve/cmd/goserve-generator@"
	getPath               = exec.LookPath
	runCmd                = exec.Command
)

func Update() {
	targetResource := fmt.Sprintf("%s%s", gitTargetInstaller, getLatestVersion())
	cmd := runCmd("go", "install", targetResource)
	_, err := cmd.CombinedOutput()

	if err != nil {
		log.Panicf("Failed to update: %v", err)
	}

	fmt.Print("✅  goserve-generator updated successfully")
	checkVersion()
}

func CheckCurrentVersion() {
	path, err := getPath(executableName)
	fmt.Println("")

	if err != nil {
		log.Panicf("Could not find goserve-generator: %v", err)
	}

	cmd := runCmd("go", "version", "-m", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Panicf("Failed to check version: %v", err)
	}

	// Parse the output to find the version
	currentVersion := extractCurrentVersion(string(output))
	if currentVersion == "" {
		fmt.Println("Could not determine version")
		return
	}

	fmt.Printf("goserve-generator version: %s\n", currentVersion)

	latestVersion := getLatestVersion()
	if latestVersion != currentVersion {
		fmt.Printf("A new version of goserve-generator is available: %s\n", latestVersion)
		fmt.Printf("goserve-generator update to get the latest version")
	}
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
