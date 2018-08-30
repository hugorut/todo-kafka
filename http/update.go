package http

import (
	"github.com/hugorut/todo-kafka/kafka"
	"net/http"
	"log"
)

type updateTodoRequest struct {
	ID int `json:"id"`
	Name string `json:"name"`
}

// updateTodoHandler exposes a json endpoint to update a Todo
// it then raises a kafka message with the provided information delegating
// responsibility for the update action
func updateTodoHandler(producer kafka.Producer) http.HandlerFunc  {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			request.Body.Close()
		}()

		if request.Method != http.MethodPut {
			http.NotFound(writer, request)
			return
		}

		var req updateTodoRequest
		if err := decodeOrError(writer, request.Body, &req); err != nil {
			return
		}

		if err := sendKafkaMessage(producer, kafka.Todos, kafka.NewUpdateTodoMessage(req.ID, req.Name)); err != nil {
			log.Printf("error encountered sending kafka msg err: %+v", err)
			respondJson(writer, http.StatusBadRequest, map[string]string{
				"error": "could not update todo",
			})
			return
		}

		respondJson(writer, http.StatusOK, map[string]string{
			"msg": "success",
		})
	}
}
