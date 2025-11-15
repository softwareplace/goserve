package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/softwareplace/goserve/cmd/goserve-generator/template"
)

type Release struct {
	TagName string `json:"tag_name"`
}

var (
	url           = "https://api.github.com/repos/softwareplace/goserve/releases/latest"
	extractBody   = getBody
	GoServeLatest = goServeLatest
)

func goServeLatest() string {
	latest, err := getLatest()
	if err != nil {
		log.Printf("Failed to fetch latest version: %v", err)
		return template.GoServeLatestVersion
	}
	return latest
}

func getLatest() (string, error) {
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

	body, err := extractBody(resp)

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

func getBody(resp *http.Response) ([]byte, error) {
	return io.ReadAll(resp.Body)
}
