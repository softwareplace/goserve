package string

import (
	"testing"
)

func TestToIntOrElseNil(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected *int
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty string",
			input:    strPtr(""),
			expected: nil,
		},
		{
			name:     "valid integer",
			input:    strPtr("42"),
			expected: intPtr(42),
		},
		{
			name:     "invalid integer",
			input:    strPtr("abc"),
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToIntOrElseNil(tt.input)
			if !compareIntPtrs(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestToIntOrElse(t *testing.T) {
	tests := []struct {
		name         string
		input        *string
		defaultValue int
		expected     *int
	}{
		{
			name:         "nil input, default used",
			input:        nil,
			defaultValue: 100,
			expected:     intPtr(100),
		},
		{
			name:         "empty string, default used",
			input:        strPtr(""),
			defaultValue: 50,
			expected:     intPtr(50),
		},
		{
			name:         "valid integer",
			input:        strPtr("25"),
			defaultValue: 75,
			expected:     intPtr(25),
		},
		{
			name:         "invalid integer, default used",
			input:        strPtr("xyz"),
			defaultValue: 300,
			expected:     intPtr(300),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToIntOrElse(tt.input, tt.defaultValue)
			if !compareIntPtrs(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func strPtr(value string) *string {
	return &value
}

func intPtr(value int) *int {
	return &value
}

func compareIntPtrs(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b != nil && *a == *b {
		return true
	}
	return false
}
