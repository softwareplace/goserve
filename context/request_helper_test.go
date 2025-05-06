package context

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockPrincipal struct {
	Id    string
	Roles []string
}

func (m *MockPrincipal) GetId() string             { return "mockId" }
func (m *MockPrincipal) GetRoles() []string        { return []string{"role1", "role2"} }
func (m *MockPrincipal) EncryptedPassword() string { return "encryptedMockPassword" }

func TestRequest_QueryOf(t *testing.T) {
	tests := []struct {
		name        string
		queryValues map[string][]string
		key         string
		expected    string
	}{
		{"key present", map[string][]string{"test": {"value1"}}, "test", "value1"},
		{"key absent", map[string][]string{"test": {"value1"}}, "missing", ""},
		{"key present with empty value", map[string][]string{"test": {""}}, "test", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request[*MockPrincipal]{QueryValues: tt.queryValues}
			result := req.QueryOf(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRequest_QueriesOf(t *testing.T) {
	tests := []struct {
		name        string
		queryValues map[string][]string
		key         string
		expected    []string
	}{
		{"key present single value", map[string][]string{"test": {"value1"}}, "test", []string{"value1"}},
		{"key absent", map[string][]string{"test": {"value1"}}, "missing", nil},
		{"key present with multiple values", map[string][]string{"test": {"value1", "value2"}}, "test", []string{"value1", "value2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request[*MockPrincipal]{QueryValues: tt.queryValues}
			result := req.QueriesOf(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRequest_QueryOfOrElse(t *testing.T) {
	tests := []struct {
		name        string
		queryValues map[string][]string
		key         string
		defaultVal  string
		expected    string
	}{
		{"key present", map[string][]string{"test": {"value1"}}, "test", "default", "value1"},
		{"key absent", map[string][]string{"test": {"value1"}}, "missing", "default", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request[*MockPrincipal]{QueryValues: tt.queryValues}
			result := req.QueryOfOrElse(tt.key, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRequest_HeadersOf(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string][]string
		key      string
		expected []string
	}{
		{"key present", map[string][]string{"test": {"value1"}}, "test", []string{"value1"}},
		{"key absent", map[string][]string{"test": {"value1"}}, "missing", nil},
		{"key present with multiple values", map[string][]string{"test": {"value1", "value2"}}, "test", []string{"value1", "value2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request[*MockPrincipal]{Headers: tt.headers}
			result := req.HeadersOf(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRequest_HeaderOfOrElse(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string][]string
		key        string
		defaultVal string
		expected   string
	}{
		{"key present", map[string][]string{"test": {"value1"}}, "test", "default", "value1"},
		{"key absent", map[string][]string{"test": {"value1"}}, "missing", "default", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request[*MockPrincipal]{Headers: tt.headers}
			result := req.HeaderOfOrElse(tt.key, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRequest_FormValue(t *testing.T) {
	r := &http.Request{
		Form: url.Values{
			"test": {"value1"},
		},
	}
	req := &Request[*MockPrincipal]{Request: r}

	assert.Equal(t, "value1", req.FormValue("test"))
	assert.Equal(t, "", req.FormValue("missing"))
}

func TestRequest_ParseMultipartForm(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader("test"))
	req := &Request[*MockPrincipal]{Request: r}

	err := req.ParseMultipartForm(32 << 20)
	assert.Equal(t, http.ErrNotMultipart, err)
}
