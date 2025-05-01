package utils

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue string
		expected     string
		setupEnv     bool
	}{
		{"env exists and not empty", "TEST_ENV", "test_value", "default_value", "test_value", true},
		{"env does not exist", "TEST_ENV", "", "default_value", "default_value", false},
		{"env exists but empty", "TEST_ENV", "", "default_value", "default_value", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnv {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			} else {
				os.Unsetenv(tt.envKey)
			}

			got := GetEnvOrDefault(tt.envKey, tt.defaultValue)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestGetBoolEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue bool
		expected     bool
		setupEnv     bool
	}{
		{"env exists and true (lowercase)", "BOOL_ENV", "true", false, true, true},
		{"env exists and true (uppercase)", "BOOL_ENV", "TRUE", false, true, true},
		{"env exists and false", "BOOL_ENV", "false", true, false, true},
		{"env exists and invalid value", "BOOL_ENV", "invalid", true, false, true},
		{"env does not exist", "BOOL_ENV", "", true, true, false},
		{"env exists but empty", "BOOL_ENV", "", false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnv {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			} else {
				os.Unsetenv(tt.envKey)
			}

			got := GetBoolEnvOrDefault(tt.envKey, tt.defaultValue)
			if got != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, got)
			}
		})
	}
}
