package http

import (
	"testing"
)

type data struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func TestServerReloading(t *testing.T) {
	config := Build("https://jsonplaceholder.typicode.com").
		WithPath("todos/1")

	t.Run("must to parse response body to string", func(t *testing.T) {
		api := NewService()

		response, err := api.Get(config)

		if err != nil {
			t.Error(err)
		}

		if response == nil {
			t.Error("Received response is empty")
		}

		str, err := api.ToString()

		if err != nil {
			t.Error(err)
		}

		t.Log(str)
	})

	t.Run("must to parse response body to expected type struct", func(t *testing.T) {
		api := NewService()

		response, err := api.Get(config)

		if err != nil {
			t.Error(err)
		}

		if response == nil {
			t.Error("Received response is empty")
		}

		responseData := data{}
		err = api.BodyDecode(&responseData)

		if err != nil {
			t.Error(err)
		}

		if responseData.UserID == 0 {
			t.Error("UserID field is not valid or missing")
		}

		if responseData.ID == 0 {
			t.Error("ID field is not valid or missing")
		}

		if responseData.Title == "" {
			t.Error("Title field is not valid or missing")
		}
	})
}
