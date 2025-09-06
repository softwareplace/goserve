package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)


// APIContextPath returns the API context path from the environment variable "CONTEXT_PATH".
// If the environment variable is not set or is empty, it returns "/".
func APIContextPath() string {
	if contextPath := os.Getenv("CONTEXT_PATH"); contextPath != "" {
		return ContextPathFix(contextPath)
	}
	return "/"
}

func ContextPathFix(contextPath string) string {
	contextPath = "/" + strings.TrimPrefix(contextPath, "/")
	return strings.TrimSuffix(contextPath, "/") + "/"
}


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

// GetBoolEnvOrDefault retrieves a boolean environment variable value or returns a default.
// It converts the environment variable string value to a boolean, considering "true" (case-insensitive) as true
// and all other values as false.
//
// Parameters:
//   - env: The name of the environment variable to retrieve
//   - defaultValue: The default boolean value to return if the environment variable is not set or empty
//
// Returns:
//   - The boolean value of the environment variable if set and not empty
//   - The defaultValue if the environment variable is not set or empty
func GetBoolEnvOrDefault(env string, defaultValue bool) bool {
	value, exists := os.LookupEnv(env)
	if !exists || value == "" {
		return defaultValue
	}
	return strings.ToLower(value) == "true"
}

// GetRequiredEnv retrieves a required environment variable value.
// If the environment variable is not set or empty, it logs a fatal error and exits the program.
//
// Parameters:
//   - key: The name of the required environment variable
//
// Returns:
//   - The value of the environment variable if set and not empty
//
// Exits with log.Fatal if the environment variable is not set or empty
func GetRequiredEnv(key string) string {
	envValue, exists := os.LookupEnv(key)
	if !exists || envValue == "" {
		log.Panic(key + " environment variable is required")
	}
	return envValue
}

// GetRequiredIntEnv retrieves a required integer environment variable value.
// If the environment variable is not set, empty, or cannot be converted to an integer,
// it logs a fatal error and exits the program.
//
// Parameters:
//   - key: The name of the required integer environment variable
//
// Returns:
//   - The integer value of the environment variable
//
// Exits with log.Fatal if the environment variable is invalid
func GetRequiredIntEnv(key string) int {
	value := GetRequiredEnv(key)
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Panicf(key+" failed with error: ", err)
	}
	return intValue
}

// GetRequiredInt64Env retrieves a required 64-bit integer environment variable value.
// If the environment variable is not set, empty, or cannot be converted to an int64,
// it logs a fatal error and exits the program.
//
// Parameters:
//   - key: The name of the required int64 environment variable
//
// Returns:
//   - The int64 value of the environment variable
//
// Exits with log.Fatal if the environment variable is invalid
func GetRequiredInt64Env(key string) int64 {
	value := GetRequiredEnv(key)
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Panicf(key+" failed with error: ", err)
	}
	return intValue
}

// GetRequiredFloat64Env retrieves a required 64-bit floating-point environment variable value.
// If the environment variable is not set, empty, or cannot be converted to a float64,
// it logs a fatal error and exits the program.
//
// Parameters:
//   - key: The name of the required float64 environment variable
//
// Returns:
//   - The float64 value of the environment variable
//
// Exits with log.Fatal if the environment variable is invalid
func GetRequiredFloat64Env(key string) float64 {
	value := GetRequiredEnv(key)
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Panicf(key+" failed with error: ", err)
	}
	return floatValue
}

// GetRequiredBoolEnv retrieves a required boolean environment variable value.
// If the environment variable is not set, empty, or cannot be converted to a boolean,
// it logs a fatal error and exits the program.
//
// Parameters:
//   - key: The name of the required boolean environment variable
//
// Returns:
//   - The boolean value of the environment variable
//
// Exits with log.Fatal if the environment variable is invalid
func GetRequiredBoolEnv(key string) bool {
	value := GetRequiredEnv(key)
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Panicf(key+" failed with error: ", err)
	}
	return boolValue
}

// GetIntEnvOrElseDefault retrieves an integer environment variable value or returns a default.
// If the environment variable is not set or cannot be converted to an integer,
// it returns the provided default value.
//
// Parameters:
//   - key: The name of the environment variable to retrieve
//   - defaultValue: The default integer value to return if the environment variable is not set or empty
func GetIntEnvOrElseDefault(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)

	if err != nil {
		log.Warnf("Failed to parse %s as int: %v", key, err)
		return defaultValue
	}

	return intValue
}
