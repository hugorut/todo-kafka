package todo

import (
	"github.com/hugorut/todo-kafka/kafka"
	"encoding/json"
	"log"
	"github.com/hugorut/todo-kafka/storage"
	"reflect"
	"strings"
)


// unmarshal takes a slice of bytes and json unmarshals it to a given interface
// handling the error if there is a problem in the conversion
func unmarshal(msg []byte, val interface{}) error {
	err := json.Unmarshal(msg, val)

	if err != nil {
		t := reflect.TypeOf(val)
		names := strings.Split(t.String(), ".")
		name := names[len(names)-1]

		log.Printf("error unmarshalling %s err: %+v", name, err)
		return err
	}

	return nil
}

// MessageCallback defines an interface which Broker calls when
// an inbound ConsumerMessage is to be handled
type MessageCallback func(message []byte)

// Create is a MessageCallback which stores a given storage.Todo in a
// the underlying memory
func Create(msg []byte) {
	var create kafka.CreateTodoMessage
	if err := unmarshal(msg, &create); err != nil {
		return
	}

	log.Printf("creating todo with name: %s", create.Name)
	storage.S.Put(storage.Todo{
		Name: create.Name,
	})
}

// Update is a MessageCallback which updates a storage.Todo
// with the given kafka message information
func Update(msg []byte) {
	var update kafka.UpdateTodoMessage
	if err := unmarshal(msg, &update); err != nil {
		return
	}

	log.Printf("update todo with id: %v name: %s", update.ID,  update.Name)
	storage.S.Put(storage.Todo{
		ID: update.ID,
		Name: update.Name,
	})
}

// Delete is a MessageCallback which removes a given Todo
// from the storage using the kafka message ID
func Delete(msg []byte) {
	var del kafka.DeleteTodoMessage
	if err := unmarshal(msg, &del); err != nil {
		return
	}

	log.Printf("deleting todo with id: %v", del.ID)
	storage.S.Del(del.ID)
}
