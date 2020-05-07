package sendy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	TestTodoObject struct {
		UserID int    `json:"userId`
		Title  string `json:"title`
	}

	TestPostObject struct {
		ID int `json:"id"`
	}
)

func TestGet(t *testing.T) {
	var response TestTodoObject

	err := Get("https://jsonplaceholder.typicode.com").
		Path("/todos/1").
		SendIt().
		JSON(&response).
		Error()

	assert.Nil(t, err)
	assert.Equal(t, 1, response.UserID)
	assert.Equal(t, "delectus aut autem", response.Title)
}

func TestPost(t *testing.T) {
	var response TestPostObject

	err := Post("https://jsonplaceholder.typicode.com").
		Path("/posts").
		JSON(map[string]int{
			"id": 101,
		}).
		SendIt().
		JSON(&response).
		Error()

	assert.Nil(t, err)
	assert.Equal(t, 101, response.ID)
}
