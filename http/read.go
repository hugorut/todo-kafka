package http

import (
	"net/http"
	"github.com/hugorut/todo-kafka/storage"
)

// listTodosHandler defines a json endpoint that returns the list of
// current stored todos
func listTodosHandler(writer http.ResponseWriter, request *http.Request) {
	respondJson(writer, http.StatusOK, storage.S.List())
}
