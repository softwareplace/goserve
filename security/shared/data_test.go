package shared

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestGetRolesForPath(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		requestPath    string
		expectedRoles  []string
		expectedExists bool
	}{
		{
			name:           "Pattern match for POST::/product/:id/add-review",
			method:         "POST",
			path:           "/user/:userId/catalogs/:catalogId/view",
			requestPath:    "/user/lswvctezynpdfnhkycyugyk/catalogs/vpchihrnkzbrzvomzytwfvd/view",
			expectedRoles:  []string{"user:catalogs:view"},
			expectedExists: true,
		},
		{
			name:           "Exact match for GET::/user/profile",
			method:         "GET",
			path:           "/user/profile",
			requestPath:    "/user/profile",
			expectedRoles:  []string{"admin", "user"},
			expectedExists: true,
		},
		{
			name:           "Exact match for POST::/user/update",
			method:         "POST",
			path:           "/user/update",
			requestPath:    "/user/update",
			expectedRoles:  []string{"admin"},
			expectedExists: true,
		},
		{
			name:           "Pattern match for GET::/product/:id",
			method:         "GET",
			path:           "/product/:productId",
			requestPath:    "/product/123",
			expectedRoles:  []string{"user"},
			expectedExists: true,
		},
		{
			name:           "Pattern match for POST::/product/:id/add-review",
			method:         "POST",
			path:           "/product/:productId/add-review",
			requestPath:    "/product/456/add-review",
			expectedRoles:  []string{"user"},
			expectedExists: true,
		},
		{
			name:           "No match for non-existent route",
			method:         "GET",
			path:           "/nonexistent/path",
			requestPath:    "/nonexistent/path",
			expectedRoles:  nil,
			expectedExists: false,
		},
	}
	for _, tt := range tests {
		path := tt.method + "::" + tt.path
		AddOpenPath(path)
		AddRoles(path, tt.expectedRoles...)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock request
			req := &http.Request{
				Method: tt.method,
				URL:    &url.URL{Path: tt.requestPath},
			}

			// Call the function
			gotRoles, gotExists := GetRolesForPath(req)

			// Compare results
			if !reflect.DeepEqual(gotRoles, tt.expectedRoles) {
				t.Errorf("expected roles %v, got %v", tt.expectedRoles, gotRoles)
			}
			if gotExists != tt.expectedExists {
				t.Errorf("expected exists %v, got %v", tt.expectedExists, gotExists)
			}
		})
	}
}
