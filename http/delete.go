package http

import (
	"github.com/hugorut/todo-kafka/kafka"
	"net/http"
	"log"
)

type deleteTodoRequest struct {
	ID int `json:"id"`
}

// deleteTodoHandler exposes a json endpoint to delete a Todo
// it then raises a kafka message with the provided information delegating
// responsibility for the write action
func deleteTodoHandler(producer kafka.Producer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			request.Body.Close()
		}()

		if request.Method != http.MethodDelete {
			http.NotFound(writer, request)
			return
		}

		var req deleteTodoRequest
		if err := decodeOrError(writer, request.Body, &req); err != nil {
			return
		}

		if err := sendKafkaMessage(producer, kafka.Todos, kafka.NewDeleteTodoMessage(req.ID)); err != nil {
			log.Printf("error encountered sending kafka msg err: %+v", err)
			respondJson(writer, http.StatusBadRequest, map[string]string{
				"error": "could not delete todo",
			})
			return
		}

		respondJson(writer, http.StatusOK, map[string]string{
			"msg": "success",
		})
	}
}

