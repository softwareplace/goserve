package cmd

import (
	"log"
	"os"
	"os/exec"
)

// Execute runs a given command with optional arguments and pipes the standard
// output and error directly to the current process outputs. It ignores any
// errors that occur during the execution of the command.
func Execute(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

// MandatoryExecute runs a given command with optional arguments and pipes the
// standard output and error directly to the current process outputs. If the
// command fails to execute, it logs the error and stops the program.
func MandatoryExecute(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ Failed to Execute command '%s %v': %v", command, args, err)
	}
}
