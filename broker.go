package main

import (
	"github.com/Shopify/sarama"
	"github.com/hugorut/todo-kafka/todo"
)

// MsgChan is shorthand type for for a ConsumerMessage channel
type MsgChan chan *sarama.ConsumerMessage


// subscription is a struct which links a MsgChan to its MessageCallback
type subscription struct {
	messages MsgChan
	callback todo.MessageCallback
}

// NewSubscription returns a unexported subscription so that it is explicit that
// when used external to the package that both messages and callback properties
// are required for the Broker to function correctly
func NewSubscription(messages MsgChan, callback todo.MessageCallback) subscription {
	return subscription{
		messages: messages,
		callback: callback,
	}
}

// Broker holds a list of subscriptions that need executing
// its job is to run the subscriptions and direct messages
// to their correct callback
type Broker struct {
	subscriptions []subscription
}

// NewBroker returns a pointer to a Broker instance initializing
// the underlying subscriptions
func NewBroker(subscriptions ...subscription) *Broker  {
	return &Broker{
		subscriptions: subscriptions,
	}
}

// Broker.Run iterates through the brokers subscriptions and fires off
// a go function listening on the channels defined under the subscriptions
func (b *Broker) Run() {
	for _, s := range b.subscriptions  {
		go b.listen(s)
	}
}

// Broker.listen executes a subscription callback when a new message is received
func (b *Broker) listen(s subscription) {
	for {
		select {
		case m := <- s.messages:
			s.callback(m.Value)
		}
	}
}

