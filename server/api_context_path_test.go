package server

import (
	"os"
	"testing"

	utils "github.com/softwareplace/goserve/utils"
)

func TestApiContextPath(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		expectedPath string
	}{
		{
			name:         "Env variable not set, should return default /",
			envValue:     "",
			expectedPath: "/",
		},
		{
			name:         "Env variable set to simple path, should return formatted path",
			envValue:     "api",
			expectedPath: "/api/",
		},
		{
			name:         "Env variable starts with /, should return formatted path",
			envValue:     "/api",
			expectedPath: "/api/",
		},
		{
			name:         "Env variable ends with /, should return formatted path",
			envValue:     "api/",
			expectedPath: "/api/",
		},
		{
			name:         "Env variable starts and ends with /, should return formatted path",
			envValue:     "/api/",
			expectedPath: "/api/",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set up environment variable
			if test.envValue != "" {
				_ = os.Setenv("CONTEXT_PATH", test.envValue)
			} else {
				_ = os.Unsetenv("CONTEXT_PATH")
			}

			// Execute the function
			result := utils.APIContextPath()

			// Assert the result
			if result != test.expectedPath {
				t.Errorf("Expected '%s', but got '%s'", test.expectedPath, result)
			}
		})
	}
}
