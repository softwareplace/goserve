package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type TestStruct struct {
	Name     string `validate:"required,min=3,max=20"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,password"`
}

func TestStructValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid input",
			input: TestStruct{
				Name:     "John",
				Email:    "john.doe@example.com",
				Password: "gDszOxF0xcq6nYeR6$&$5",
			},
			wantErr: false,
		},
		{
			name: "Valid max 20 characters",
			input: TestStruct{
				Name:     "John has a long name that exceeds the max length of 20 characters",
				Email:    "john.doe@example.com",
				Password: "gDszOxF0xcq6nYeR6$&$5",
			},
			wantErr: true,
			errMsg:  "Name must be at most 20 characters",
		},
		{
			name: "Valid exactly 20 characters long name",
			input: TestStruct{
				Name:     "John has a long name",
				Email:    "john.doe@example.com",
				Password: "gDszOxF0xcq6nYeR6$&$5",
			},
			wantErr: false,
		},
		{
			name: "Missing required fields",
			input: TestStruct{
				Name:     "",
				Email:    "",
				Password: "",
			},
			wantErr: true,
			errMsg:  "Name is a required field\nEmail is a required field\nPassword is a required field",
		},
		{
			name: "Invalid email format",
			input: TestStruct{
				Name:     "John",
				Email:    "invalid-email",
				Password: "gDszOxF0xcq6nYeR6$&$5",
			},
			wantErr: true,
			errMsg:  "Email must be a valid email address",
		},
		{
			name: "Name too short",
			input: TestStruct{
				Name:     "Jo",
				Email:    "john.doe@example.com",
				Password: "gDszOxF0xcq6nYeR6$&$5",
			},
			wantErr: true,
			errMsg:  "Name must be at least 3 characters",
		},
		{
			name: "Invalid password format",
			input: TestStruct{
				Name:     "John",
				Email:    "john.doe@example.com",
				Password: "password",
			},
			wantErr: true,
			errMsg:  "Password must contain at least: 8 characters, 1 uppercase, 1 lowercase, 1 number, and 1 special character",
		},
		{
			name: "Multiple validation errors",
			input: TestStruct{
				Name:     "",
				Email:    "invalid-email",
				Password: "pass",
			},
			wantErr: true,
			errMsg:  "Name is a required field\nEmail must be a valid email address\nPassword must contain at least: 8 characters, 1 uppercase, 1 lowercase, 1 number, and 1 special character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StructValidation(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
