package http

import (
	"testing"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"github.com/hugorut/todo-kafka/storage"
	"encoding/json"
)

func Test_listTodosHandler(t *testing.T) {
	id := storage.S.Put(storage.Todo{
		Name: "test",
	})
	id2 := storage.S.Put(storage.Todo{
		Name: "test2",
	})
	r := httptest.NewRequest("GET", "/todo/list", nil)
	w := httptest.NewRecorder()

	listTodosHandler(w, r)

	assert.Equal(t,200, w.Code)

	var todos []storage.Todo
	json.Unmarshal(w.Body.Bytes(), &todos)

	assert.Contains(t, todos, storage.Todo{
		ID: id,
		Name: "test",
	})
	assert.Contains(t, todos, storage.Todo{
		ID: id2,
		Name: "test2",
	})

	storage.S.Flush()
}
