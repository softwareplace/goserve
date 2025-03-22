package utils

import (
	"os"
	"strings"
)

// GetEnvOrDefault retrieves the value of an environment variable.
// If the environment variable is not set or is empty, it returns the provided default value.
//
// Parameters:
//   - env: The name of the environment variable to retrieve.
//   - defaultValue: The default value to return if the environment variable is not set or is empty.
//
// Returns:
//   - The value of the environment variable if it is set and not empty.
//   - The defaultValue if the environment variable is not set or is empty.
//
// Example Usage:
//
//	// Load the "API_KEY" environment variable, defaulting to "defaultKey123" if not set.
//	apiKey := GetEnvOrDefault("API_KEY", "defaultKey123")
//	fmt.Println("API Key:", apiKey)
func GetEnvOrDefault(env string, defaultValue string) string {
	value, exists := os.LookupEnv(env)
	if !exists || value == "" {
		return defaultValue
	}
	return value
}

func GetBoolEnvOrDefault(env string, defaultValue bool) bool {
	value, exists := os.LookupEnv(env)
	if !exists || value == "" {
		return defaultValue
	}
	return strings.ToLower(value) == "true"
}
