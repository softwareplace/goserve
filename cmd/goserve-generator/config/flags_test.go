package config

import (
	"flag"
	"github.com/softwareplace/goserve/cmd/goserve-generator/version"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitFlagsWithSet(t *testing.T) {
	t.Cleanup(func() {
		// Reset global vars after each test
		ProjectName = ""
		Username = ""
		ReplaceCurrent = ""
		GiInit = ""
		CodeGenConfigFile = ""
		GoServerVersion = ""
	})

	args := []string{
		"-n", "myproject",
		"-u", "myuser",
		"-r", "true",
		"-gi", "false",
		"-cgf", "template.yml",
		"-gsv", "v1.2.3",
	}

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	err := InitFlagsWithSet(fs, args)
	assert.NoError(t, err)

	assert.Equal(t, "myproject", ProjectName)
	assert.Equal(t, "myuser", Username)
	assert.Equal(t, "true", ReplaceCurrent)
	assert.Equal(t, "false", GiInit)
	assert.Equal(t, "template.yml", CodeGenConfigFile)
	assert.Equal(t, "v1.2.3", GoServerVersion)
}

func TestInitFlagsWithSet_MissingRequiredFlags(t *testing.T) {
	called := false
	oldExit := osExit
	osExit = func(code int) {
		called = true
	}
	defer func() { osExit = oldExit }()

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	args := []string{}
	_ = InitFlagsWithSet(fs, args)

	assert.True(t, called, "Expected osExit to be called due to missing required flags")
}

func TestInitFlagsWithSet_InvalidFlags(t *testing.T) {
	defer func() { osExit = os.Exit }()
	exitCode := -1
	osExit = func(code int) {
		exitCode = code
	}

	InitFlags()
	require.Equal(t, 1, exitCode)
}

func TestInitFlagsWithSet_UtilitiesArgs(t *testing.T) {
	testSource := []struct {
		name     string
		args     []string
		call     func()
		expected string
		status   int
	}{
		{
			name: "should exit with status 0 for update arg",
			args: []string{"update"},
			call: func() {
				update = func() {

				}
			},
			expected: "update",
			status:   0,
		},
		{
			name: "should exit with status 0 for version arg",
			args: []string{"version"},
			call: func() {
				checkVersion = func() {

				}
			},
			expected: "version",
			status:   0,
		},
	}

	for _, tt := range testSource {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				// Reset global vars after each test
				ProjectName = ""
				Username = ""
				ReplaceCurrent = ""
				GiInit = ""
				CodeGenConfigFile = ""
				GoServerVersion = ""
				osExit = os.Exit
				checkVersion = version.CheckCurrentVersion
				update = version.Update
			})

			exitCode := -1

			osExit = func(code int) {
				exitCode = code
			}

			tt.call()

			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			err := InitFlagsWithSet(fs, tt.args)
			assert.NoError(t, err)

			require.Equal(t, tt.status, exitCode)
		})
	}
}

func TestInitFlagsWithSet_InvalidArgs(t *testing.T) {
	t.Run("should exit with status 1 for invalid args", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		args := []string{"-invalid", "arg"}
		err := InitFlagsWithSet(fs, args)
		require.Error(t, err)
	})
}
