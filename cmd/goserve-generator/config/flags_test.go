package config

import (
	"flag"
	"github.com/softwareplace/goserve/cmd/goserve-generator/version"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVersionChecker is a mock for version package
type MockVersionChecker struct {
	mock.Mock
}

func (m *MockVersionChecker) CheckCurrentVersion() {
	m.Called()
}

func (m *MockVersionChecker) Update() {
	m.Called()
}

// Test variables to restore original implementations
var (
	originalArgs          = os.Args
	originalExit          = osExit
	originalVersionCheck  = version.CheckCurrentVersion
	originalVersionUpdate = version.Update
)

func restoreOriginals() {
	os.Args = originalArgs
	osExit = originalExit
	checkVersion = originalVersionCheck
	update = originalVersionUpdate
}

func TestInitFlags(t *testing.T) {
	t.Cleanup(restoreOriginals)
	t.Cleanup(func() { flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) })

	test := struct {
		name        string
		args        []string
		expectPanic bool
		expected    struct {
			ProjectName    string
			Username       string
			ReplaceCurrent string
			GiInit         string
		}
	}{
		name: "all flags",
		args: []string{"cmd", "-n", "myproject", "-u", "myuser"},
		expected: struct {
			ProjectName    string
			Username       string
			ReplaceCurrent string
			GiInit         string
		}{
			ProjectName:    "myproject",
			Username:       "myuser",
			ReplaceCurrent: "false",
			GiInit:         "true",
		},
	}

	t.Run(test.name, func(t *testing.T) {

		os.Args = test.args

		InitFlags()

		assert.Equal(t, test.expected.ProjectName, ProjectName)
		assert.Equal(t, test.expected.Username, Username)
		assert.Equal(t, test.expected.ReplaceCurrent, ReplaceCurrent)
		assert.Equal(t, test.expected.GiInit, GiInit)
		//assert.Equal(t, flag.Usage, flagUsage)
	})
}

func TestInitFlags2(t *testing.T) {
	t.Cleanup(restoreOriginals)
	t.Cleanup(func() { flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) })

	test := struct {
		name        string
		args        []string
		expectPanic bool
		expected    struct {
			ProjectName    string
			Username       string
			ReplaceCurrent string
			GiInit         string
		}
	}{
		name: "all flags",
		args: []string{"cmd", "-n", "myproject", "-u", "myuser", "-r", "true", "-gi", "false"},
		expected: struct {
			ProjectName    string
			Username       string
			ReplaceCurrent string
			GiInit         string
		}{
			ProjectName:    "myproject",
			Username:       "myuser",
			ReplaceCurrent: "true",
			GiInit:         "false",
		},
	}

	t.Run(test.name, func(t *testing.T) {

		os.Args = test.args

		InitFlags()

		assert.Equal(t, test.expected.ProjectName, ProjectName)
		assert.Equal(t, test.expected.Username, Username)
		assert.Equal(t, test.expected.ReplaceCurrent, ReplaceCurrent)
		assert.Equal(t, test.expected.GiInit, GiInit)
	})
}

func TestArgsValidation(t *testing.T) {
	t.Cleanup(restoreOriginals)
	t.Cleanup(func() { flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) })

	// Mock os.Exit
	exitCalled := false
	exitCode := 0
	osExit = func(code int) {
		exitCalled = true
		exitCode = code
	}

	// Mock version functions
	mockVersion := new(MockVersionChecker)
	checkVersion = mockVersion.CheckCurrentVersion
	update = mockVersion.Update

	tests := []struct {
		name           string
		args           []string
		setupFlags     func()
		expectExit     bool
		expectExitCode int
		expectVersion  bool
		expectUpdate   bool
	}{
		{
			name:           "version command",
			args:           []string{"cmd", "version"},
			expectExit:     true,
			expectExitCode: 0,
			expectVersion:  true,
		},
		{
			name:           "update command",
			args:           []string{"cmd", "update"},
			expectExit:     true,
			expectExitCode: 0,
			expectUpdate:   true,
		},
		{
			name: "missing required flags",
			args: []string{"cmd"},
			setupFlags: func() {
				ProjectName = ""
				Username = ""
			},
			expectExit:     true,
			expectExitCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exitCalled = false
			os.Args = tt.args

			if tt.setupFlags != nil {
				tt.setupFlags()
			}

			if tt.expectVersion {
				mockVersion.On("CheckCurrentVersion").Once()
			}
			if tt.expectUpdate {
				mockVersion.On("Update").Once()
			}

			argsValidation()

			assert.Equal(t, tt.expectExit, exitCalled)
			if tt.expectExit {
				assert.Equal(t, tt.expectExitCode, exitCode)
			}

			mockVersion.AssertExpectations(t)
		})
	}
}

func TestInit(t *testing.T) {
	// Verify log flags are set correctly
	assert.Equal(t, 0, log.Flags())

}
