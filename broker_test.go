package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/Shopify/sarama"
	"time"
	"github.com/hugorut/todo-kafka/kafka"
	"github.com/hugorut/todo-kafka/todo"
	"fmt"
)

func TestBroker_Run_CallsCorrectSubscription(t *testing.T) {
	c := make(MsgChan, 1)
	var called bool
	callback := func(b []byte) {
		called = true
	}

	c2 := make(MsgChan, 1)
	var called2 bool
	callback2 := func(b []byte) {
		called2 = true
	}

	b := NewBroker(NewSubscription(c, map[kafka.Event]todo.MessageCallback{
		kafka.CreateTodo: callback,
	}), NewSubscription(c2, map[kafka.Event]todo.MessageCallback{
		kafka.CreateTodo: callback2,
	}))
	b.Run()

	c2 <- &sarama.ConsumerMessage{
		Value: []byte(fmt.Sprintf(`{"Type":"%s"}`, kafka.CreateTodo)),
	}
	time.Sleep(1 * time.Second)
	assert.False(t, called)
	assert.True(t, called2)
}
