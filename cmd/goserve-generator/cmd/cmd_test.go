package cmd

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name          string
		command       string
		args          []string
		expectedError bool
	}{
		{"validCommand", "echo", []string{"Hello, World!"}, false},
		{"invalidCommand", "nonexistent_command", []string{}, true},
		{"validCommandWithArgs", "echo", []string{"Test", "123"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create buffers to capture stdout and stderr
			stdout := &bytes.Buffer{}

			// Override os.Stdout and os.Stderr temporarily
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			defer func() {
				os.Stdout = oldStdout
				os.Stderr = oldStderr
			}()
			r, w, _ := os.Pipe()
			os.Stdout = w
			os.Stderr = w

			Execute(tt.command, tt.args...)

			w.Close()
			os.Stdout = oldStdout
			os.Stderr = oldStderr

			// Read from the pipe
			stdout.ReadFrom(r)

			// Check for valid/invalid command cases
			if tt.expectedError {
				if stdout.Len() > 0 {
					t.Errorf("expected no output for invalid command, got: %s", stdout.String())
				}
			} else {
				if stdout.Len() == 0 {
					t.Errorf("expected output, got none")
				}
			}
		})
	}
}

func TestMandatoryExecute(t *testing.T) {
	tests := []struct {
		name          string
		command       string
		args          []string
		expectedPanic bool
	}{
		{"validCommand", "echo", []string{"Hello, Mandatory Test!"}, false},
		{"invalidCommand", "nonexistent_command", []string{}, true},
		{"validCommandWithArgs", "echo", []string{"Another", "Test"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture logs
			var logOutput bytes.Buffer
			log.SetOutput(&logOutput)
			defer func() {
				log.SetOutput(os.Stderr)
			}()

			// Capture panic if expected
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectedPanic {
						t.Errorf("unexpected panic: %v", r)
					}
				} else if tt.expectedPanic {
					t.Errorf("expected panic, but did not occur")
				}
			}()

			MandatoryExecute(tt.command, tt.args...)

			// Check for valid command logging
			if !tt.expectedPanic && tt.command == "echo" && logOutput.Len() > 0 {
				t.Errorf("unexpected log output: %s", logOutput.String())
			}
		})
	}
}
