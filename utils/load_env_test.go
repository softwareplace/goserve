package utils

import (
	goserveerror "github.com/softwareplace/goserve/error"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// testFatal is a helper function to test if a function calls log.Fatal
func testFatal(t *testing.T, fn func()) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()
	fn()
	return false
}

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		setEnv       bool
		defaultValue string
		want         string
	}{
		{
			name:         "unset env returns default",
			envKey:       "TEST_STR_1",
			setEnv:       false,
			defaultValue: "default",
			want:         "default",
		},
		{
			name:         "empty env returns default",
			envKey:       "TEST_STR_2",
			envValue:     "",
			setEnv:       true,
			defaultValue: "default",
			want:         "default",
		},
		{
			name:         "set env returns value",
			envKey:       "TEST_STR_3",
			envValue:     "custom_value",
			setEnv:       true,
			defaultValue: "default",
			want:         "custom_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}
			got := GetEnvOrDefault(tt.envKey, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetEnvOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBoolEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		setEnv       bool
		defaultValue bool
		want         bool
	}{
		{
			name:         "unset env returns default true",
			envKey:       "TEST_BOOL_1",
			setEnv:       false,
			defaultValue: true,
			want:         true,
		},
		{
			name:         "empty env returns default false",
			envKey:       "TEST_BOOL_2",
			envValue:     "",
			setEnv:       true,
			defaultValue: false,
			want:         false,
		},
		{
			name:         "true value returns true",
			envKey:       "TEST_BOOL_3",
			envValue:     "true",
			setEnv:       true,
			defaultValue: false,
			want:         true,
		},
		{
			name:         "TRUE value returns true",
			envKey:       "TEST_BOOL_4",
			envValue:     "TRUE",
			setEnv:       true,
			defaultValue: false,
			want:         true,
		},
		{
			name:         "horse value returns false",
			envKey:       "TEST_BOOL_5",
			envValue:     "horse",
			setEnv:       true,
			defaultValue: true,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}
			got := GetBoolEnvOrDefault(tt.envKey, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetBoolEnvOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequiredEnv(t *testing.T) {
	// Test successful case
	t.Run("existing env returns value", func(t *testing.T) {
		expected := "test-value"
		os.Setenv("TEST_REQUIRED", expected)
		defer os.Unsetenv("TEST_REQUIRED")

		got := GetRequiredEnv("TEST_REQUIRED")
		if got != expected {
			t.Errorf("GetRequiredEnv() = %v, want %v", got, expected)
		}
	})

	// Test error cases
	t.Run("missing env triggers fatal", func(t *testing.T) {
		os.Setenv("NONEXISTENT_ENV", "")
		defer os.Unsetenv("NONEXISTENT_ENV")

		failed := false

		goserveerror.Handler(func() {
			GetRequiredEnv("NONEXISTENT_ENV")
		}, func(err error) {
			failed = true
		})

		require.True(t, failed)
	})

	t.Run("empty env triggers fatal", func(t *testing.T) {
		os.Setenv("EMPTY_ENV", "")
		defer os.Unsetenv("EMPTY_ENV")

		failed := false

		goserveerror.Handler(func() {
			GetRequiredEnv("EMPTY_ENV")
		}, func(err error) {
			failed = true
		})

		require.True(t, failed)
	})
}

func TestGetRequiredIntEnv(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		envValue string
		want     int
		wantErr  bool
	}{
		{
			name:     "valid integer",
			envKey:   "TEST_INT",
			envValue: "123",
			want:     123,
			wantErr:  false,
		},
		{
			name:     "invalid integer (horse)",
			envKey:   "TEST_INT",
			envValue: "horse",
			wantErr:  true,
		},
		{
			name:     "empty value",
			envKey:   "TEST_INT",
			envValue: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.envKey, tt.envValue)
			defer os.Unsetenv(tt.envKey)

			if tt.wantErr {
				if !testFatal(t, func() {
					GetRequiredIntEnv(tt.envKey)
				}) {
					t.Error("Expected fatal error")
				}
				return
			}

			got := GetRequiredIntEnv(tt.envKey)
			if got != tt.want {
				t.Errorf("GetRequiredIntEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequiredInt64Env(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		envValue string
		want     int64
		wantErr  bool
	}{
		{
			name:     "valid int64",
			envKey:   "TEST_INT64",
			envValue: "9223372036854775807",
			want:     9223372036854775807,
			wantErr:  false,
		},
		{
			name:     "invalid int64 (horse)",
			envKey:   "TEST_INT64",
			envValue: "horse",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.envKey, tt.envValue)
			defer os.Unsetenv(tt.envKey)

			if tt.wantErr {
				if !testFatal(t, func() {
					GetRequiredInt64Env(tt.envKey)
				}) {
					t.Error("Expected fatal error")
				}
				return
			}

			got := GetRequiredInt64Env(tt.envKey)
			if got != tt.want {
				t.Errorf("GetRequiredInt64Env() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequiredFloat64Env(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		envValue string
		want     float64
		wantErr  bool
	}{
		{
			name:     "valid float",
			envKey:   "TEST_FLOAT",
			envValue: "123.456",
			want:     123.456,
			wantErr:  false,
		},
		{
			name:     "invalid float (horse)",
			envKey:   "TEST_FLOAT",
			envValue: "horse",
			wantErr:  true,
		},
		{
			name:     "scientific notation",
			envKey:   "TEST_FLOAT",
			envValue: "1.23e-4",
			want:     0.000123,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.envKey, tt.envValue)
			defer os.Unsetenv(tt.envKey)

			if tt.wantErr {
				if !testFatal(t, func() {
					GetRequiredFloat64Env(tt.envKey)
				}) {
					t.Error("Expected fatal error")
				}
				return
			}

			got := GetRequiredFloat64Env(tt.envKey)
			if got != tt.want {
				t.Errorf("GetRequiredFloat64Env() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequiredBoolEnv(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		envValue string
		want     bool
		wantErr  bool
	}{
		{
			name:     "true value",
			envKey:   "TEST_BOOL",
			envValue: "true",
			want:     true,
			wantErr:  false,
		},
		{
			name:     "invalid bool (horse)",
			envKey:   "TEST_BOOL",
			envValue: "horse",
			wantErr:  true,
		},
		{
			name:     "1 value",
			envKey:   "TEST_BOOL",
			envValue: "1",
			want:     true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.envKey, tt.envValue)
			defer os.Unsetenv(tt.envKey)

			if tt.wantErr {
				if !testFatal(t, func() {
					GetRequiredBoolEnv(tt.envKey)
				}) {
					t.Error("Expected fatal error")
				}
				return
			}

			got := GetRequiredBoolEnv(tt.envKey)
			if got != tt.want {
				t.Errorf("GetRequiredBoolEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIntEnvOrElseDefault(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		setEnv       bool
		defaultValue int
		want         int
	}{
		{
			name:         "unset env returns default",
			envKey:       "TEST_INT_1",
			setEnv:       false,
			defaultValue: 42,
			want:         42,
		},
		{
			name:         "empty env returns default",
			envKey:       "TEST_INT_2",
			envValue:     "",
			setEnv:       true,
			defaultValue: 99,
			want:         99,
		},
		{
			name:         "valid int env returns value",
			envKey:       "TEST_INT_3",
			envValue:     "123",
			setEnv:       true,
			defaultValue: 7,
			want:         123,
		},
		{
			name:         "invalid int env returns default",
			envKey:       "TEST_INT_4",
			envValue:     "notanint",
			setEnv:       true,
			defaultValue: 55,
			want:         55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			} else {
				os.Unsetenv(tt.envKey)
			}
			got := GetIntEnvOrElseDefault(tt.envKey, tt.defaultValue)
			require.Equal(t, tt.want, got)
		})
	}
}
