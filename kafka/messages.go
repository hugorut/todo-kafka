package kafka

import (
	"github.com/Shopify/sarama"
	"encoding/json"
)

type CreateTodoMessage struct {
	Name string
}

type DeleteTodoMessage struct {
	ID int
}

type UpdateTodoMessage struct {
	ID int
	Name string
}

// NewMessage provides a simple way to new up a ProducerMessage which holds
// marshaled json in the Value property
func NewMessage(topic Topic, partition int32, msg interface{}) *sarama.ProducerMessage {
	b,_ := json.Marshal(msg)

	return &sarama.ProducerMessage{
		Topic: topic.String(),
		Partition: partition,
		Value: sarama.ByteEncoder(b),
	}
}
