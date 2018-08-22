package http

import (
	"testing"
	"encoding/json"
	"net/http/httptest"
	"bytes"
	"github.com/Shopify/sarama/mocks"
	"github.com/hugorut/todo-kafka/kafka"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/hugorut/todo-kafka/test"
	"log"
	"errors"
)

func Test_deleteTodoHandler_SendsKafkaMessage(t *testing.T) {
	b, _ := json.Marshal(deleteTodoRequest{
		ID: 15,
	})

	r := httptest.NewRequest("DELETE", "/todo/delete", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	p := mocks.NewSyncProducer(t, nil)
	p.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {
		var m kafka.DeleteTodoMessage
		json.Unmarshal(val, &m)

		if 15 !=  m.ID {
			return fmt.Errorf("expecting DeleteTodoMessge to hold id: %v instead received %v", 15 , m.ID)
		}

		return nil
	})

	deleteTodoHandler(p)(w, r)
}

func Test_deleteTodoHandler_respondsSuccess(t *testing.T) {
	b, _ := json.Marshal(deleteTodoRequest{
		ID: 15,
	})
	r := httptest.NewRequest("DELETE", "/todo/delete", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	p := mocks.NewSyncProducer(t, nil)
	p.ExpectSendMessageAndSucceed()

	deleteTodoHandler(p)(w, r)

	assert.Equal(t,200, w.Code)
	assert.JSONEq(t, `{"msg":"success"}`, w.Body.String())
}

func Test_deleteTodoHandler_respondsErrorIfProducerFail(t *testing.T) {
	writer := &test.Writer{}
	log.SetOutput(writer)

	b, _ := json.Marshal(createTodoRequest{
		Name: "test",
	})
	r := httptest.NewRequest("DELETE", "/todo/delete", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	p := mocks.NewSyncProducer(t, nil)
	p.ExpectSendMessageAndFail(errors.New("my message"))

	deleteTodoHandler(p)(w, r)

	assert.Equal(t,400, w.Code)
	assert.JSONEq(t, `{"error":"could not delete todo"}`, w.Body.String())
	assert.Contains(t, writer.Msg(), "error encountered sending kafka msg err: my message")
}
