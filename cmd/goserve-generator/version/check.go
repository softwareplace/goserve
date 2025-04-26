package version

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

func Update() {
	targetResource := fmt.Sprintf("github.com/softwareplace/goserve/cmd/goserve-generator@%s", GoServeLatest())
	cmd := exec.Command("go", "install", targetResource)
	_, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatalf("Failed to update: %v", err)
	}

	fmt.Print("✅  goserve-generator updated successfully")
	CheckCurrentVersion()
}

func CheckCurrentVersion() {
	path, err := exec.LookPath("goserve-generator")
	fmt.Println("")

	if err != nil {
		log.Fatalf("Could not find goserve-generator: %v", err)
	}

	cmd := exec.Command("go", "version", "-m", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to check version: %v", err)
	}

	// Parse the output to find the version
	currentVersion := extractVersion(string(output))
	if currentVersion == "" {
		fmt.Println("Could not determine version")
		return
	}

	fmt.Printf("goserve-generator version: %s\n", currentVersion)

	latestVersion := GoServeLatest()
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
