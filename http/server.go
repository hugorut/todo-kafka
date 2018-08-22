package http

import (
	"net/http"
	"github.com/hugorut/todo-kafka/kafka"
	"encoding/json"
	"log"
	"io"
	"fmt"
)

// respondJson provides a quick function to write json to the ResponseWriter
// setting the appropriate headers before doing so
func respondJson(w http.ResponseWriter, status int, msg interface{})  {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// decodeOrError wraps commonly used decoding functionality. If the decode fails
// it writes to a log and then responds to the client with a fitting error message
func decodeOrError(writer http.ResponseWriter, reader io.Reader, req interface{}) error {
	err := json.NewDecoder(reader).Decode(req)
	if err != nil {
		log.Printf("error unmarshalling json request to todo handler: %+v", err)
		respondJson(writer, http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})

		return err
	}

	return nil
}

// sendKafkaMessage provides a more fluent interface to the kafka.SendMessage function
func sendKafkaMessage(producer kafka.Producer, topic kafka.Topic, msg interface{}) error  {
	return kafka.SendMessage(
		producer,
		kafka.NewMessage(
			topic,
			0,
			msg,
		),
	)
}

// NewServer sets up a http server at the give port & initializes given todo handlers
func NewServer(producer kafka.Producer, port string) error  {
	http.HandleFunc("/todo/create", createTodoHandler(producer))
	http.HandleFunc("/todo/update", updateTodoHandler(producer))
	http.HandleFunc("/todo/delete", deleteTodoHandler(producer))
	http.HandleFunc("/todo/list", listTodosHandler)

	return http.ListenAndServe(fmt.Sprintf(":%s",port), nil)
}
