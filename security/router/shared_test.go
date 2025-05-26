package router

import (
	"testing"
)

func TestAddOpenPath(t *testing.T) {
	tests := []struct {
		name          string
		inputPath     string
		existingPaths []string
		expectedPaths []string
	}{
		{"add new path", "/new/path", []string{}, []string{"/new/path"}},
		{"add duplicate path", "/new/path", []string{"/new/path"}, []string{"/new/path"}},
		{"normalize path", "//path//to//resource", []string{}, []string{"/path/to/resource"}},
		{"normalize and avoid duplicate", "//path//to//resource", []string{"/path/to/resource"}, []string{"/path/to/resource"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			openPathLock.Lock()
			openPaths = append([]string{}, tt.existingPaths...)
			openPathLock.Unlock()

			AddOpenPath(tt.inputPath)

			openPathLock.Lock()
			defer openPathLock.Unlock()

			if len(openPaths) != len(tt.expectedPaths) {
				t.Fatalf("expected %v paths, got %v", len(tt.expectedPaths), len(openPaths))
			}
			for i, path := range openPaths {
				if path != tt.expectedPaths[i] {
					t.Fatalf("expected path %v, got %v", tt.expectedPaths[i], path)
				}
			}
		})
	}
}

func TestAddRoles(t *testing.T) {
	tests := []struct {
		name          string
		inputPath     string
		inputRoles    []string
		existingRoles map[string][]string
		expectedRoles map[string][]string
	}{
		{"add roles to new path", "/path1", []string{"role1", "role2"}, map[string][]string{}, map[string][]string{"/path1": {"role1", "role2"}}},
		{"overwrite roles for existing path", "/path1", []string{"role3"}, map[string][]string{"/path1": {"role1"}}, map[string][]string{"/path1": {"role3"}}},
		{"do not add if no roles provided", "/path2", []string{}, map[string][]string{}, map[string][]string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roles = make(map[string][]string)
			for k, v := range tt.existingRoles {
				roles[k] = v
			}

			AddRoles(tt.inputPath, tt.inputRoles...)

			if len(roles) != len(tt.expectedRoles) {
				t.Fatalf("expected %v roles, got %v", len(tt.expectedRoles), len(roles))
			}
			for path, expectedRoles := range tt.expectedRoles {
				actualRoles, exists := roles[path]
				if !exists {
					t.Fatalf("expected roles for path %v, but none found", path)
				}
				if len(actualRoles) != len(expectedRoles) {
					t.Fatalf("expected %v roles for path %v, got %v", len(expectedRoles), path, len(actualRoles))
				}
				for i, role := range actualRoles {
					if role != expectedRoles[i] {
						t.Fatalf("expected role %v at index %v for path %v, got %v", expectedRoles[i], i, path, role)
					}
				}
			}
		})
	}
}
