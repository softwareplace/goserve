package login

import (
	"github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security/encryptor"
	"testing"
)

// MockPrincipal is a mock implementation of goservectx.Principal interface
type MockPrincipal struct {
	encryptedPassword string
}

func (m *MockPrincipal) GetId() string {
	return "mock-id"
}

func (m *MockPrincipal) GetRoles() []string {
	return []string{"role1", "role2"}
}

func (m *MockPrincipal) EncryptedPassword() string {
	return m.encryptedPassword
}

func TestDefaultPasswordValidator_IsValidPassword(t *testing.T) {
	encrypt := encryptor.NewEncrypt("valid_password").EncodedPassword()

	tests := []struct {
		name           string
		loginEntryData User
		principal      context.Principal
		expected       bool
	}{
		{
			name: "valid password",
			loginEntryData: User{
				Password: "valid_password",
			},
			principal: &MockPrincipal{
				encryptedPassword: encrypt,
			},
			expected: true,
		},
		{
			name: "invalid password",
			loginEntryData: User{
				Password: "invalid_password",
			},
			principal: &MockPrincipal{
				encryptedPassword: encrypt,
			},
			expected: false,
		},
		{
			name: "empty password input",
			loginEntryData: User{
				Password: "",
			},
			principal: &MockPrincipal{
				encryptedPassword: encrypt,
			},
			expected: false,
		},
		{
			name: "empty encrypted password",
			loginEntryData: User{
				Password: "valid_password",
			},
			principal: &MockPrincipal{
				encryptedPassword: "",
			},
			expected: false,
		},
		{
			name: "both empty passwords",
			loginEntryData: User{
				Password: "",
			},
			principal: &MockPrincipal{
				encryptedPassword: "",
			},
			expected: false,
		},
	}

	validator := &DefaultPasswordValidator[context.Principal]{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.IsValidPassword(tt.loginEntryData, tt.principal)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
