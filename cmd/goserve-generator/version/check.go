package version

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

func CheckCurrentVersion() {
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
	currentVersion := extractVersion(string(output))
	if currentVersion == "" {
		fmt.Println("Could not determine version")
		return
	}

	fmt.Printf("goserve-generator version: %s\n", currentVersion)
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
