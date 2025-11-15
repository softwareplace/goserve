package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	Name     string  `validate:"required,min=3,max=20"`
	Email    string  `validate:"required,email"`
	Password string  `validate:"required,password"`
	Age      int     `validate:"gte=1,lte=150"`
	Amount   float32 `validate:"gt=10,lte=1000"`
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
				Age:      100,
				Amount:   100,
			},
			wantErr: false,
		},
		{
			name: "Valid max 20 characters",
			input: TestStruct{
				Name:     "John has a long name that exceeds the max length of 20 characters",
				Email:    "john.doe@example.com",
				Password: "gDszOxF0xcq6nYeR6$&$5",
				Age:      100,
				Amount:   100,
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
				Age:      100,
				Amount:   100,
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
			errMsg:  "Name is a required field\nEmail is a required field\nPassword is a required field\nAge must be greater or equal to 1\nAmount must be greater than 10",
		},
		{
			name: "Invalid email format",
			input: TestStruct{
				Name:     "John",
				Email:    "invalid-email",
				Password: "gDszOxF0xcq6nYeR6$&$5",
				Age:      100,
				Amount:   100,
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
				Age:      100,
				Amount:   100,
			},
			wantErr: true,
			errMsg:  "Name must be at least 3 characters",
		},
		{
			name: "Invalid min number",
			input: TestStruct{
				Name:     "John has a long name",
				Email:    "john.doe@example.com",
				Password: "gDszOxF0xcq6nYeR6$&$5",
				Age:      0,
				Amount:   0,
			},
			wantErr: true,
			errMsg:  "Age must be greater or equal to 1\nAmount must be greater than 10",
		},
		{
			name: "Invalid max number",
			input: TestStruct{
				Name:     "John has a long name",
				Email:    "john.doe@example.com",
				Password: "gDszOxF0xcq6nYeR6$&$5",
				Age:      1000,
				Amount:   1001,
			},
			wantErr: true,
			errMsg:  "Age must be less or equal to 150\nAmount must be less or equal to 1000",
		},
		{
			name: "Valid number",
			input: TestStruct{
				Name:     "John has a long name",
				Email:    "john.doe@example.com",
				Password: "gDszOxF0xcq6nYeR6$&$5",
				Age:      100,
				Amount:   100,
			},
			wantErr: false,
		},
		{
			name: "Invalid password format",
			input: TestStruct{
				Name:     "John",
				Email:    "john.doe@example.com",
				Password: "password",
				Age:      100,
				Amount:   100,
			},
			wantErr: true,
			errMsg:  "Password must contain at least: 8 characters, 1 uppercase, 1 lowercase, 1 number, and 1 special character",
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
