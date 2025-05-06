package context

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type mockPrincipal struct{}

func (m *mockPrincipal) GetId() string             { return "mockId" }
func (m *mockPrincipal) GetRoles() []string        { return []string{"role1", "role2"} }
func (m *mockPrincipal) EncryptedPassword() string { return "encryptedMockPassword" }

func newMockContext() *Request[*mockPrincipal] {
	req := httptest.NewRequest(http.MethodGet, "/path?queryKey=queryValue", nil)
	req.Header.Set("X-Api-Key", "mockApiKey")
	req.Header.Set("Authorization", "mockAuth")
	res := httptest.NewRecorder()
	return Of[*mockPrincipal](res, req, "testReference")
}

func TestRequest_Context(t *testing.T) {
	ctx := newMockContext()
	assert.Equal(t, "mockApiKey", ctx.ApiKey)
	assert.Equal(t, "mockAuth", ctx.Authorization)
}

func TestRequest_GetSample(t *testing.T) {
	ctx := newMockContext()
	sample := ctx.GetSample()
	assert.Equal(t, ctx.ApiKey, sample.ApiKey)
}

func TestRequest_Flush(t *testing.T) {
	ctx := newMockContext()
	ctx.Flush()
	assert.Nil(t, ctx.Writer)
	assert.Nil(t, ctx.Request)
	assert.Nil(t, ctx.Principal)
}

func TestRequest_HeaderOf(t *testing.T) {
	ctx := newMockContext()
	assert.Equal(t, "mockApiKey", ctx.HeaderOf("X-Api-Key"))
	assert.Equal(t, "default", ctx.HeaderOfOrElse("Non-Existent", "default"))
}

func TestRequest_QueryMethods(t *testing.T) {
	ctx := newMockContext()
	assert.Equal(t, "queryValue", ctx.QueryOf("queryKey"))
	assert.Empty(t, ctx.QueryOf("nonKey"))
	assert.Equal(t, []string{"queryValue"}, ctx.QueriesOf("queryKey"))
	assert.Empty(t, ctx.QueriesOf("nonKey"))
	assert.Equal(t, []string{"default"}, ctx.QueriesOfElse("key", []string{"default"}))
}

func TestRequest_PathValueOf(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test/paramValue", nil)
	res := httptest.NewRecorder()
	m := mux.NewRouter()
	m.HandleFunc("/test/{param}", func(w http.ResponseWriter, r *http.Request) {
		ctx := Of[*mockPrincipal](w, r, "testReference")
		assert.Equal(t, "paramValue", ctx.PathValueOf("param"))
	}).Methods(http.MethodGet)

	req = mux.SetURLVars(req, map[string]string{"param": "paramValue"})
	m.ServeHTTP(res, req)
}

func TestRequest_Write(t *testing.T) {
	ctx := newMockContext()
	body := map[string]string{"key": "value"}
	ctx.Write(body, http.StatusOK)
	result := httptest.NewRecorder()
	json.NewEncoder(result).Encode(body)
	assert.JSONEq(t, result.Body.String(), (*ctx.Writer).(*httptest.ResponseRecorder).Body.String())
}

func TestRequest_ErrorResponses(t *testing.T) {
	expectedError := []struct {
		error    string
		callback func(ctx *Request[*mockPrincipal])
	}{
		{
			error: "Unauthorized",
			callback: func(ctx *Request[*mockPrincipal]) {
				ctx.Unauthorized()
			},
		},
		{
			error: "Bad Request",
			callback: func(ctx *Request[*mockPrincipal]) {
				ctx.BadRequest("Bad Request")
			},
		},
		{
			error: "Internal server error",
			callback: func(ctx *Request[*mockPrincipal]) {
				ctx.InternalServerError("Internal server error")
			},
		},
		{
			error: "Not Found",
			callback: func(ctx *Request[*mockPrincipal]) {
				ctx.NotFount("Not Found")
			},
		},
		{
			error: "No content",
			callback: func(ctx *Request[*mockPrincipal]) {
				ctx.NoContent("No content")
			},
		},
		{
			error: "Created",
			callback: func(ctx *Request[*mockPrincipal]) {
				ctx.Created("Created")
			},
		},
		{
			error: "Forbidden",
			callback: func(ctx *Request[*mockPrincipal]) {
				ctx.Forbidden("Forbidden")
			},
		},
		{
			error: "Invalid input",
			callback: func(ctx *Request[*mockPrincipal]) {
				ctx.InvalidInput()
			},
		},
	}

	for _, err := range expectedError {
		ctx := newMockContext()
		err.callback(ctx)
		bodyString := (*ctx.Writer).(*httptest.ResponseRecorder).Body.String()

		if strings.Contains(bodyString, err.error) {
			t.Logf("Response body is good %s", bodyString)
		} else {
			t.Errorf("Expected response body to contain %s on %s", err.error, bodyString)
		}
	}
}

func TestRequest_WriteReader(t *testing.T) {
	ctx := newMockContext()
	content := []byte("fileContent")
	reader := bytes.NewReader(content)
	err := ctx.WriteReader(reader, "test.txt")
	assert.NoError(t, err)
	assert.Contains(t, (*ctx.Writer).(*httptest.ResponseRecorder).Header().Get("Content-Disposition"), "attachment")
}

func TestRequest_CreateNew(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/path", nil)
	res := httptest.NewRecorder()
	ctx := Of[*mockPrincipal](res, req, "reference")
	assert.NotNil(t, ctx)
	assert.NotEmpty(t, ctx.sessionId)
}

func TestRequest_FormFile(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	formFile, _ := writer.CreateFormFile("file", "test.txt")
	formFile.Write([]byte("test data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res := httptest.NewRecorder()
	ctx := Of[*mockPrincipal](res, req, "testReference")

	file, _, err := ctx.FormFile("file")
	assert.NoError(t, err)
	content := make([]byte, 9)
	file.Read(content)
	assert.Equal(t, "test data", string(content))
}
