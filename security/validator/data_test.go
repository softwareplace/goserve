package validator

import (
	"github.com/softwareplace/http-utils/api_context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
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
			method:         "POST",
			path:           "/user/:userId/catalogs/:catalogId/view",
			requestPath:    "/user/lswvctezynpdfnhkycyugyk/catalogs/vpchihrnkzbrzvomzytwfvd/view",
			expectedRoles:  []string{"user:catalogs:view"},
			expectedExists: true,
		},
		{
			method:         "GET",
			path:           "/user/profile",
			requestPath:    "/user/profile",
			expectedRoles:  []string{"admin", "user"},
			expectedExists: true,
		},
		{
			method:         "POST",
			path:           "/user/update",
			requestPath:    "/user/update",
			expectedRoles:  []string{"admin"},
			expectedExists: true,
		},
		{
			method:         "GET",
			path:           "/product/:productId",
			requestPath:    "/product/123",
			expectedRoles:  []string{"user"},
			expectedExists: true,
		},
		{
			method:         "POST",
			path:           "/product/:productId/add-review",
			requestPath:    "/product/456/add-review",
			expectedRoles:  []string{"user"},
			expectedExists: true,
		},
		{
			method:         "GET",
			path:           "/nonexistent/path",
			requestPath:    "/nonexistent/path",
			expectedRoles:  nil,
			expectedExists: false,
		},
	}
	for _, tt := range tests {
		path := tt.method + "::" + tt.path
		AddRoles(path, tt.expectedRoles...)
	}

	for _, tt := range tests {
		t.Run("given__"+tt.path+"==>"+tt.requestPath+"__must_return__"+strconv.FormatBool(tt.expectedExists), func(t *testing.T) {

			// Create a mock request

			ctx := api_context.Of[*api_context.DefaultContext](httptest.NewRecorder(), &http.Request{
				Method: tt.method,
				URL:    &url.URL{Path: tt.requestPath},
			}, "")

			// Call the function
			gotRoles, gotExists := GetRolesForPath(*ctx)

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

// Test for IsPublicPath method
func TestForPublicPaths(t *testing.T) {
	t.Run("Test IsPublicPath implementation", func(t *testing.T) {
		tests := []struct {
			name           string
			method         string
			path           string
			requestPath    string
			expectedResult bool
		}{
			{
				method:         "GET",
				path:           "/public/info",
				requestPath:    "/public/info",
				expectedResult: true,
			},
			{
				method:         "POST",
				path:           "/user/restricted-access",
				requestPath:    "/user/restricted-access",
				expectedResult: false,
			},
			{
				method:         "GET",
				path:           "/api/:version/doc",
				requestPath:    "/api/v2/doc",
				expectedResult: true,
			},
			{
				method:         "GET",
				path:           "/admin/:id",
				requestPath:    "/admin/158",
				expectedResult: false,
			},
		}
		for _, tt := range tests {
			path := tt.method + "::" + tt.path
			if tt.expectedResult {
				AddOpenPath(path) // Add the path to open/public paths
			}
		}

		for _, tt := range tests {
			t.Run("given__"+tt.path+"==>"+tt.requestPath+"__must_return__"+strconv.FormatBool(tt.expectedResult), func(t *testing.T) {

				// Create mock request
				ctx := api_context.Of[*api_context.DefaultContext](httptest.NewRecorder(), &http.Request{
					Method: tt.method,
					URL:    &url.URL{Path: tt.requestPath},
				}, "")
				// Call the function
				isPublic := IsPublicPath(*ctx)
				// Compare results
				if isPublic != tt.expectedResult {
					t.Errorf("expected public %v, got %v", tt.expectedResult, isPublic)
				}
			})
		}
	})
}

// Test for HasResourceAccessRight function
func TestHasResourceAccessRight(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		requestPath    string
		userRoles      []string
		requiredRoles  []string
		isRoleRequired bool
		expectedResult bool
	}{
		{
			name:           "User with matching role has access",
			method:         "GET",
			path:           "/admin/panel",
			requestPath:    "/admin/panel",
			userRoles:      []string{"admin"},
			requiredRoles:  []string{"admin"},
			isRoleRequired: true,
			expectedResult: true,
		},
		{
			name:           "User without matching role does not have access",
			method:         "POST",
			path:           "/user/create",
			requestPath:    "/user/create",
			userRoles:      []string{"user"},
			requiredRoles:  []string{"admin"},
			isRoleRequired: true,
			expectedResult: false,
		},
		{
			name:           "Path does not require roles, access granted",
			method:         "GET",
			path:           "/public/info",
			requestPath:    "/public/info",
			userRoles:      []string{},
			requiredRoles:  nil,
			isRoleRequired: false,
			expectedResult: true,
		},
		{
			name:           "Path requires roles but user has none",
			method:         "GET",
			path:           "/restricted/data",
			requestPath:    "/restricted/data",
			userRoles:      []string{},
			requiredRoles:  []string{"user"},
			isRoleRequired: true,
			expectedResult: false,
		},
		{
			name:           "User has one of the required roles",
			method:         "GET",
			path:           "/validator/resource",
			requestPath:    "/validator/resource",
			userRoles:      []string{"user", "editor"},
			requiredRoles:  []string{"admin", "editor"},
			isRoleRequired: true,
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add roles for the path
			path := tt.method + "::" + tt.path
			if tt.isRoleRequired && len(tt.requiredRoles) > 0 {
				AddRoles(path, tt.requiredRoles...)
			}

			// Create mock context
			ctx := api_context.Of[*api_context.DefaultContext](httptest.NewRecorder(), &http.Request{
				Method: tt.method,
				URL:    &url.URL{Path: tt.requestPath},
			}, "")

			ctx.RequestData = &api_context.DefaultContext{}

			ctx.RequestData.SetRoles(tt.userRoles)

			rolesLoader := func(ctx api_context.ApiRequestContext[*api_context.DefaultContext]) []string {
				return ctx.RequestData.Roles()
			}

			ctx.AccessRolesLoader = &rolesLoader

			// Call the function
			result := HasResourceAccessRight(*ctx)

			// Compare results
			if result != tt.expectedResult {
				t.Errorf("expected result %v, got %v", tt.expectedResult, result)
			}
		})
	}
}
