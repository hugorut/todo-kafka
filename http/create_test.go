package http

import (
	"testing"
	"net/http/httptest"
	"bytes"
	"encoding/json"
	"github.com/Shopify/sarama/mocks"
	"github.com/hugorut/todo-kafka/kafka"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/pkg/errors"
	"github.com/hugorut/todo-kafka/test"
	"log"
)

func Test_createTodoHandler_SendsKafkaMessage(t *testing.T) {
	b, _ := json.Marshal(createTodoRequest{
		Name: "test",
	})
	r := httptest.NewRequest("POST", "/todo/create", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	p := mocks.NewSyncProducer(t, nil)
	p.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {
		var m kafka.CreateTodoMessage
		json.Unmarshal(val, &m)

		if "test" !=  m.Name {
			return fmt.Errorf("expecting CreateTodoMessage to hold name: %s instead received %s", "test", m.Name)
		}

		return nil
	})

	createTodoHandler(p)(w, r)
}

func Test_createTodoHandler_respondsSuccess(t *testing.T) {
	b, _ := json.Marshal(createTodoRequest{
		Name: "test",
	})
	r := httptest.NewRequest("POST", "/todo/create", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	p := mocks.NewSyncProducer(t, nil)
	p.ExpectSendMessageAndSucceed()

	createTodoHandler(p)(w, r)

	assert.Equal(t,200, w.Code)
	assert.JSONEq(t, `{"msg":"success"}`, w.Body.String())
}

func Test_createTodoHandler_respondsErrorIfProducerFail(t *testing.T) {
	writer := &test.Writer{}
	log.SetOutput(writer)

	b, _ := json.Marshal(createTodoRequest{
		Name: "test",
	})
	r := httptest.NewRequest("POST", "/todo/create", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	p := mocks.NewSyncProducer(t, nil)
	p.ExpectSendMessageAndFail(errors.New("my message"))

	createTodoHandler(p)(w, r)

	assert.Equal(t,400, w.Code)
	assert.JSONEq(t, `{"error":"could not process todo"}`, w.Body.String())
	assert.Contains(t, writer.Msg(), "error encountered sending kafka msg err: my message")
}
