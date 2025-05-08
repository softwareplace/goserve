package request

import (
	"github.com/gorilla/mux"
	"github.com/softwareplace/goserve/context"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockFormBody struct {
	DirName  string   `json:"dirName"`
	FileName string   `json:"fileName"`
	Tags     []string `json:"tags"`
}

type MockRequest struct {
	XApiKey string       `name:"X-Api-Key" required:"true" validate:"required" error_message:"required header param [X-Api-Key]" header:"X-Api-Key" json:"X-Api-Key"`
	Page    int          `name:"page" required:"true" validate:"required" error_message:"required query param [page]" query:"page" json:"page"`
	Count   int          `name:"count" required:"true" validate:"required" error_message:"required query param [count]" query:"count" json:"count"`
	UserId  int          `name:"userId" required:"true" validate:"required" error_message:"required path param [userId]" path:"userId" json:"userId"`
	Body    MockFormBody `name:"body" json:"body" required:"true" validate:"required"`
}

func TestRequest_BindRequestParams(t *testing.T) {
	t.Run("should return no error when required header was provided", func(t *testing.T) {
		router := mux.NewRouter()

		var ctx *context.Request[*context.DefaultContext]
		request := MockRequest{
			Body: MockFormBody{},
		}
		var errBind *RequestError

		router.HandleFunc("/login/{userId}", func(w http.ResponseWriter, r *http.Request) {
			ctx = context.Of[*context.DefaultContext](w, r, "test")
			errBind = BindRequestParams(r, &request)
			if errBind != nil {
				ctx.InternalServerError(errBind.Error())
				return
			}
			ctx.Ok(request)
		}).Methods("POST")

		req, err := http.NewRequest("POST", "/login/101?count=1000&page=10", nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", context.MultipartFormData)
		req.Header.Set("X-Api-Key", "test")
		req.Form = make(map[string][]string)

		req.Form["dirName"] = []string{"app-test"}
		req.Form["fileName"] = []string{"app.txt"}
		req.Form["tags"] = []string{"app-test", "app-test-2"}

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code)

		require.Nil(t, errBind)
		require.Equal(t, "test", request.XApiKey)
		require.Equal(t, 10, request.Page)
		require.Equal(t, 1000, request.Count)
		require.Equal(t, 101, request.UserId)
		require.Equal(t, "app-test", request.Body.DirName)
		require.Equal(t, "app.txt", request.Body.FileName)
		require.Equal(t, []string{"app-test", "app-test-2"}, request.Body.Tags)
	})

	t.Run("should return error when required header was not provided", func(t *testing.T) {
		router := mux.NewRouter()

		var ctx *context.Request[*context.DefaultContext]
		request := MockRequest{}
		var errBind *RequestError

		router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
			ctx = context.Of[*context.DefaultContext](w, r, "test")
			errBind = BindRequestParams(r, &request)
			if errBind != nil {
				ctx.InternalServerError(errBind.Error())
				return
			}
			ctx.Ok(request)
		}).Methods("POST")

		req, err := http.NewRequest("POST", "/login", nil)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusInternalServerError, recorder.Code)

		expected := []string{
			"Key: required header param [X-Api-Key] Error:Field validation for required header param [X-Api-Key] failed on the required tag",
			"Key: required query param [page] Error:Field validation for required query param [page] failed on the required tag",
			"Key: required query param [count] Error:Field validation for required query param [count] failed on the required tag",
			"Key: required path param [userId] Error:Field validation for required path param [userId] failed on the required tag",
		}

		for _, e := range expected {
			require.True(t, strings.Contains(errBind.Error(), e), e)
		}
	})
}
