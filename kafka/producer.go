package kafka

import (
	"github.com/Shopify/sarama"
	"log"
)

// Producer is a wrapper interface for the underlying sarama.SyncProducer
type Producer interface {
	sarama.SyncProducer
}

// NewProducer wraps sarama.NewProducer changing the signature to a more readable
// variadic argument list and setting some safe defaults
func NewProducer(brokers ...string) (Producer, error) {
	config := sarama.NewConfig()

	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.ClientID = "producer"

	return sarama.NewSyncProducer(brokers, config)
}


// SendMessage wraps Producer.SendMessage with some loggin and returning a
// more manageable single error argument
func SendMessage(producer Producer, msg *sarama.ProducerMessage) error {
	p, o, err := producer.SendMessage(msg)
	log.Printf("sent message to partition: %v offset: %v\n", p, o)

	return err
}
