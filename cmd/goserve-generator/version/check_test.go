package version

import (
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCommandExecutor is a mock for command execution
type MockCommandExecutor struct {
	mock.Mock
}

func (m *MockCommandExecutor) LookPath(file string) (string, error) {
	args := m.Called(file)
	return args.String(0), args.Error(1)
}

func (m *MockCommandExecutor) Command(name string, arg ...string) ([]byte, error) {
	args := m.Called(name, arg)
	return args.Get(0).([]byte), args.Error(1)
}

// Test variables to restore original implementations
var (
	originalLookPath         = exec.LookPath
	originalCommand          = exec.Command
	originalCheckVersion     = checkVersion
	originalGetLatestVersion = getLatestVersion
)

func restoreOriginals() {
	getPath = originalLookPath
	runCmd = originalCommand
	checkVersion = originalCheckVersion
	getLatestVersion = originalGetLatestVersion
}

func TestUpdate(t *testing.T) {
	t.Cleanup(restoreOriginals)

	tests := []struct {
		name                   string
		latestVersion          string
		commandOutput          []byte
		commandError           error
		wantErr                bool
		failedOnCombinedOutput bool
	}{
		{
			name:          "successful update",
			latestVersion: "v1.2.3",
			commandOutput: []byte("success"),
			commandError:  nil,
			wantErr:       false,
		},
		{
			name:          "failed update",
			latestVersion: "v1.2.3",
			commandOutput: []byte("error"),
			commandError:  errors.New("command failed"),
			wantErr:       true,
		},
		{
			name:                   "failed combined output",
			latestVersion:          "v1.2.3",
			commandOutput:          []byte("error"),
			commandError:           errors.New("command failed"),
			wantErr:                true,
			failedOnCombinedOutput: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock getLatestVersion
			getLatestVersion = func() string {
				return tt.latestVersion
			}

			checkVersion = func() {
				if tt.wantErr && !tt.failedOnCombinedOutput {
					log.Panicf("version check failed %v", tt.commandError)
				}
			}

			// Mock command execution
			runCmd = func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "go", name)
				assert.Equal(t, []string{"install", "github.com/softwareplace/goserve/cmd/goserve-generator@" + tt.latestVersion}, arg)
				command := exec.Command("ls")
				if tt.failedOnCombinedOutput {
					_, err := command.CombinedOutput()
					require.NoError(t, err)
				}
				return command
			}

			if tt.wantErr {
				assert.Panics(t, func() { Update() })
			} else {
				assert.NotPanics(t, func() { Update() })
			}
		})
	}
}

func TestCheckCurrentVersion(t *testing.T) {
	t.Cleanup(restoreOriginals)

	tests := []struct {
		name                   string
		lookPathError          error
		commandOutput          []byte
		commandError           error
		latestVersion          string
		currentVersion         string
		wantPanic              bool
		expectOutput           string
		expectNewUpdate        bool
		failedGetPath          bool
		failedOnCombinedOutput bool
	}{
		{
			name:           "successful version check with update available",
			lookPathError:  nil,
			commandOutput:  []byte("mod\tgithub.com/softwareplace/goserve\tv1.0.0\n"),
			commandError:   nil,
			latestVersion:  "v1.1.0",
			currentVersion: "v1.0.0",
			wantPanic:      false,
			expectOutput:   "A new version of goserve-generator is available: v1.1.0",
		},
		{
			name:           "successful version check with no update",
			lookPathError:  nil,
			commandOutput:  []byte("mod\tgithub.com/softwareplace/goserve\tv1.1.0\n"),
			commandError:   nil,
			latestVersion:  "v1.1.0",
			currentVersion: "v1.1.0",
			wantPanic:      false,
			expectOutput:   "goserve-generator version: v1.1.0",
		},
		{
			name:          "executable not found",
			lookPathError: errors.New("not found"),
			wantPanic:     true,
		},
		{
			name:          "version command failed",
			lookPathError: nil,
			commandError:  errors.New("command failed"),
			wantPanic:     true,
		},
		{
			name:          "get path failed",
			lookPathError: nil,
			commandError:  errors.New("command failed"),
			failedGetPath: true,
			wantPanic:     true,
		},
		{
			name:                   "failed on combined output",
			lookPathError:          nil,
			commandError:           errors.New("command failed"),
			wantPanic:              true,
			failedOnCombinedOutput: true,
		},
		{
			name:           "version extraction failed",
			lookPathError:  nil,
			commandOutput:  []byte("invalid output"),
			commandError:   nil,
			wantPanic:      false,
			expectOutput:   "Could not determine version",
			currentVersion: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			old := log.StandardLogger().Out
			defer func() { log.SetOutput(old) }()
			var buf strings.Builder
			log.SetOutput(&buf)

			// Mock lookPath
			getPath = func(path string) (string, error) {
				assert.Equal(t, executableName, path)
				if tt.failedGetPath {
					return "", tt.commandError
				}
				if tt.wantPanic && !tt.failedOnCombinedOutput {
					log.Panicf("version check failed %v", tt.commandError)
				} else {
					log.Info(tt.latestVersion)
					if tt.expectOutput != "" {
						log.Info(tt.expectOutput)
					}
				}
				return "/path/to/goserve-generator", tt.lookPathError
			}

			// Mock command execution
			runCmd = func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "go", name)
				assert.Equal(t, []string{"version", "-m", "/path/to/goserve-generator"}, arg)
				command := exec.Command("ls")

				if tt.latestVersion != "" {
					command = exec.Command("echo", "mod\tgithub.com/softwareplace/goserve "+tt.currentVersion)
				}

				if tt.failedOnCombinedOutput {
					_, err := command.CombinedOutput()
					require.NoError(t, err)
				}
				return command
			}

			// Mock getLatestVersion
			getLatestVersion = func() string {
				return tt.latestVersion
			}

			if tt.wantPanic {
				assert.Panics(t, CheckCurrentVersion)
			} else {
				assert.NotPanics(t, CheckCurrentVersion)
				output := buf.String()
				assert.Contains(t, output, tt.expectOutput)

				if tt.expectNewUpdate {
					assert.Contains(t, output, tt.latestVersion)
				}
			}
		})
	}
}

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "mod version",
			input:    "mod\tgithub.com/softwareplace/goserve\tv1.2.3\n",
			expected: "v1.2.3",
		},
		{
			name:     "dep version",
			input:    "dep\tgithub.com/softwareplace/goserve\tv1.2.4\n",
			expected: "v1.2.4",
		},
		{
			name:     "multiple lines",
			input:    "some line\nmod\tgithub.com/softwareplace/goserve\tv1.2.5\nanother line",
			expected: "v1.2.5",
		},
		{
			name:     "no version",
			input:    "some text without version",
			expected: "",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, extractVersion(tt.input))
		})
	}
}
