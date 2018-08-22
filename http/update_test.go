package http

import (
	"testing"
	"encoding/json"
	"net/http/httptest"
	"bytes"
	"github.com/Shopify/sarama/mocks"
	"github.com/hugorut/todo-kafka/kafka"
	"fmt"
)

func Test_udpateTodoHandler_SendsKafkaMessage(t *testing.T) {
	b, _ := json.Marshal(updateTodoRequest{
		ID: 15,
		Name: "test",
	})

	r := httptest.NewRequest("PUT", "/todo/update", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	p := mocks.NewSyncProducer(t, nil)
	p.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {
		var m kafka.UpdateTodoMessage
		json.Unmarshal(val, &m)

		if 15 !=  m.ID {
			return fmt.Errorf("expecting UpdateTodoMessage to hold id: 15 instead received %v", m.ID)
		}

		if "test" !=  m.Name {
			return fmt.Errorf("expecting UpdateTodoMessage to hold name: %s instead received %s", "test", m.Name)
		}

		return nil
	})

	updateTodoHandler(p)(w, r)
}
