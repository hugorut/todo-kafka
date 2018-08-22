package kafka

import (
	"github.com/Shopify/sarama"
	"log"
)

// Consumer is a wrapper interface for the underlying sarama.Consumer
type Consumer interface {
	sarama.Consumer
}

// NewConsumer wraps sarama.NewConsumer changing the signature to a more readable
// variadic argument list
func NewConsumer(addrs ...string) (Consumer, error) {
	config := sarama.NewConfig()
	config.ClientID = "consumer"

	return sarama.NewConsumer(addrs, config)
}

// ConsumeTopic subscribes a given Consumer to a Topic. Returning a message channel
// than aggregates all of the partition messages into one stream
func ConsumeTopic(consumer Consumer, topic Topic) (chan *sarama.ConsumerMessage, error) {
	msgs := make(chan *sarama.ConsumerMessage)

	partitions, err := consumer.Partitions(topic.String())
	if err != nil {
		return msgs, err
	}

	for _, p := range partitions {
		log.Printf("consuming partition: %v for topic: %s", p, topic)
		pc, err := consumer.ConsumePartition(topic.String(), p, sarama.OffsetNewest)
		if err != nil {
			return msgs, err
		}

		go receiveMessage(pc, msgs)
	}

	return msgs, nil
}

// receiveMessage defines a blocking function that listens for inbound PartitionConsumer messages
func receiveMessage(pc sarama.PartitionConsumer, msgs chan *sarama.ConsumerMessage)  {
	for {
		select {
		case message := <-pc.Messages():
			msgs <- message
		case err := <-pc.Errors():
			log.Printf("kafka consumer error received err: %+v", err)
		}
	}
}
