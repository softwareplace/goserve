package path

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserHomePathFix(t *testing.T) {
	dir, err := os.UserHomeDir()
	require.NoError(t, err)
	tests := []struct {
		name      string
		inputPath string
		mockHome  string
		mockError error
		wantPath  string
		wantPanic bool
	}{
		{
			name:      "no_tilde",
			inputPath: "/usr/local/bin",
			mockHome:  dir,
			wantPath:  "/usr/local/bin",
		},
		{
			name:      "tilde_replacement",
			inputPath: "~/documents",
			mockHome:  dir,
			wantPath:  dir + "/documents",
		},
		{
			name:      "tilde_only",
			inputPath: "~",
			mockHome:  dir,
			wantPath:  dir,
		},
		{
			name:      "empty_input",
			inputPath: "",
			mockHome:  dir,
			wantPath:  "",
		},
		{
			name:      "error_in_user_home",
			inputPath: "~/documents",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock os.UserHomeDir

			home := os.Getenv("HOME")
			userProfile := os.Getenv("USERPROFILE")
			userProfilePercent := os.Getenv("%userprofile%")
			homeDrive := os.Getenv("HOMEDRIVE")

			if tt.wantPanic {
				_ = os.Unsetenv("HOME")
				_ = os.Unsetenv("USERPROFILE")
				_ = os.Unsetenv("%userprofile%")
				_ = os.Unsetenv("HOMEDRIVE")
			}

			defer func() {
				_ = os.Setenv("HOME", home)
				_ = os.Setenv("USERPROFILE", userProfile)
				_ = os.Setenv("%userprofile%", userProfilePercent)
				_ = os.Setenv("HOMEDRIVE", homeDrive)
			}()

			var got string

			got, err = UserHomePathFix(tt.inputPath)

			if tt.wantPanic {
				if err == nil {
					t.Error("expected panic, got nil")
				}
			} else {
				if got != tt.wantPath {
					t.Errorf("UserHomePathFix(%q) = %q, want %q", tt.inputPath, got, tt.wantPath)
				}
			}
		})
	}
}
