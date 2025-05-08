package config

import (
	"flag"
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
