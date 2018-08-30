package kafka

import (
	"github.com/Shopify/sarama"
	"encoding/json"
	"log"
)

type Event string

const (
	CreateTodo Event = "CreateTodo"
	UpdateTodo Event = "UpdateTodo"
	DeleteTodo Event = "DeleteTodo"
)

type Message struct {
	ID int
	Type Event
}

type CreateTodoMessage struct {
	Message
	Name string
}

func NewCreateTodoMessage(name string) UpdateTodoMessage {
	return UpdateTodoMessage{
		Message: Message{
			Type: CreateTodo,
		},
		Name: name,
	}
}

type DeleteTodoMessage struct {
	Message
}

func NewDeleteTodoMessage(id int) UpdateTodoMessage {
	return UpdateTodoMessage{
		Message: Message{
			ID: id,
			Type: DeleteTodo,
		},
	}
}

type UpdateTodoMessage struct {
	Message
	Name string
}

func NewUpdateTodoMessage(id int, name string) UpdateTodoMessage {
	return UpdateTodoMessage{
		Message: Message{
			ID: id,
			Type: UpdateTodo,
		},
		Name: name,
	}
}


// NewMessage provides a simple way to new up a ProducerMessage which holds
// marshaled json in the Value property
func NewMessage(topic Topic, partition int32, msg interface{}) *sarama.ProducerMessage {
	b,_ := json.Marshal(msg)
	log.Printf("creating new kafka message with val: %s", b)

	return &sarama.ProducerMessage{
		Topic: topic.String(),
		Partition: partition,
		Value: sarama.ByteEncoder(b),
	}
}
