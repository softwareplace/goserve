package version

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/softwareplace/goserve/cmd/goserve-generator/template"

	"github.com/stretchr/testify/assert"
)

func TestGoServeLatest(t *testing.T) {
	t.Run("should return default version when request failed", func(t *testing.T) {
		originalUrl := url
		url = "http://dev.invalid.url/goserve/releases/latest"
		defer func() {
			url = originalUrl
		}()

		latest := GoServeLatest()
		require.Equal(t, template.GoServeLatestVersion, latest)
	})

	t.Run("should return default version when request return 404 status code", func(t *testing.T) {
		originalUrl := url
		url = "https://api.github.com/repos/softwareplace/goserve/releases/latestversion"
		defer func() {
			url = originalUrl
		}()

		latest := GoServeLatest()
		require.Equal(t, template.GoServeLatestVersion, latest)
	})

	t.Run("should return default version when failed to extract body", func(t *testing.T) {

		extractBody = func(resp *http.Response) ([]byte, error) {
			return nil, fmt.Errorf("failed to extract body")
		}
		defer func() {
			extractBody = getBody
		}()

		latest := GoServeLatest()
		require.Equal(t, template.GoServeLatestVersion, latest)
	})

	t.Run("should return default version when tagName is not defined", func(t *testing.T) {

		extractBody = func(resp *http.Response) ([]byte, error) {
			return []byte(`{}`), nil
		}
		defer func() {
			extractBody = getBody
		}()

		latest := GoServeLatest()
		require.Equal(t, template.GoServeLatestVersion, latest)
	})
}

func TestGoServeLatest_Success(t *testing.T) {
	// Mock HTTP server response
	mockResponse := `{"tag_name": "v1.2.3"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Override the global URL variable for testing
	url = server.URL

	latest := GoServeLatest()
	assert.Equal(t, "v1.2.3", latest, "GoServeLatest should return the latest release version")
}

func TestGoServeLatest_Fallback(t *testing.T) {
	// Mock HTTP server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Override the global URL variable for testing
	url = server.URL

	latest := GoServeLatest()
	assert.Equal(t, template.GoServeLatestVersion, latest, "GoServeLatest should return the fallback version when the HTTP request fails")
}

func TestGetLatest_Success(t *testing.T) {
	// Mock HTTP server response
	mockResponse := `{"tag_name": "v1.2.3"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Override the global URL variable for testing
	url = server.URL

	latest, err := getLatest()
	assert.NoError(t, err, "getLatest should not return an error for a valid response")
	assert.Equal(t, "v1.2.3", latest, "getLatest should return the correct tag name")
}

func TestGetLatest_ErrorResponses(t *testing.T) {
	t.Run("HTTP error", func(t *testing.T) {
		// Mock HTTP server that returns an error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		// Override the global URL for testing
		url = server.URL

		latest, err := getLatest()
		assert.Error(t, err, "getLatest should return an error when HTTP response status is not OK")
		assert.Empty(t, latest, "getLatest should return an empty string on HTTP error")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		// Mock HTTP server with malformed JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`invalid json`))
		}))
		defer server.Close()

		// Override the global URL for testing
		url = server.URL

		latest, err := getLatest()
		assert.Error(t, err, fmt.Errorf("failed to parse JSON"), "getLatest should return an error when JSON parsing fails")
		assert.Empty(t, latest, "getLatest should return an empty string on JSON parsing error")
	})
}
