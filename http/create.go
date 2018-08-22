package http

import (
	"github.com/hugorut/todo-kafka/kafka"
	"net/http"
	"log"
)

type createTodoRequest struct {
	Name string `json:"name"`
}

// createTodoHandler exposes a json endpoint to create a new Todo
// it then raises a kafka message with the provided information delegating
// responsibility for the write action
func createTodoHandler(producer kafka.Producer) http.HandlerFunc  {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			request.Body.Close()
		}()

		var req createTodoRequest
		if err := decodeOrError(writer, request.Body, &req); err != nil {
			return
		}

		if err := sendKafkaMessage(producer, kafka.CreateToDo, kafka.CreateTodoMessage{ Name: req.Name }); err != nil {
			log.Printf("error encountered sending kafka msg err: %+v", err)
			respondJson(writer, http.StatusBadRequest, map[string]string{
				"error": "could not process todo",
			})
			return
		}

		respondJson(writer, http.StatusOK, map[string]string{
			"msg": "success",
		})
	}
}
