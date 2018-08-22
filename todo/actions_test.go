package todo

import (
	"testing"
	"encoding/json"
	"github.com/hugorut/todo-kafka/kafka"
	"github.com/hugorut/todo-kafka/storage"
	"github.com/stretchr/testify/assert"
	"github.com/hugorut/todo-kafka/test"
	"log"
)

func TestCreate_CallsSPut(t *testing.T) {
	b, _ := json.Marshal(kafka.CreateTodoMessage{Name: "test"})
	Create(b)

	todos := storage.S.List()
	assert.Len(t, todos, 1)
	assert.Equal(t, "test", todos[0].Name)

	storage.S.Flush()
}


func TestCreate_LogsError(t *testing.T) {
	w := &test.Writer{}
	log.SetOutput(w)

	Create([]byte(`[[asfd`))

	assert.Contains(t, w.Msg(), "error unmarshalling CreateTodoMessage err: invalid character 'a'")
}

func TestDelete_RemovesFromStorage(t *testing.T) {
	storage.S.Put(storage.Todo{
		Name: "test",
	})

	id2 := storage.S.Put(storage.Todo{
		Name: "test2",
	})

	storage.S.Put(storage.Todo{
		Name: "test3",
	})

	b, _ := json.Marshal(kafka.DeleteTodoMessage{ID: id2})
	Delete(b)

	todos := storage.S.List()
	assert.Len(t, todos, 2)

	storage.S.Flush()
}

func TestUpdate_CallsSPutWithID(t *testing.T) {
	id := storage.S.Put(storage.Todo{
		Name: "test",
	})

	b, _ := json.Marshal(kafka.UpdateTodoMessage{ID: id, Name: "test2"})
	Update(b)

	todos := storage.S.List()
	assert.Len(t, todos, 1)
	assert.Equal(t, "test2", todos[0].Name)

	storage.S.Flush()
}
