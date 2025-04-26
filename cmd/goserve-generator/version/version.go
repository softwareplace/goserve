package version

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/softwareplace/goserve/cmd/goserve-generator/template"
	"io"
	"net/http"
)

type Release struct {
	TagName string `json:"tag_name"`
}

func GoServeLatest() string {
	latest, err := getLatest()
	if err != nil {
		log.Printf("Failed to fetch latest version: %v", err)
		return template.GoServeLatestVersion
	}
	return latest
}

func getLatest() (string, error) {
	url := "https://api.github.com/repos/softwareplace/goserve/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch releases: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}

	if release.TagName == "" {
		return template.GoServeLatestVersion, nil
	}
	return release.TagName, nil
}
