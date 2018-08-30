package main

import (
	"testing"
	"github.com/Shopify/sarama"
	"time"
	"github.com/hugorut/todo-kafka/kafka"
	"github.com/hugorut/todo-kafka/todo"
	"fmt"
)

func TestBroker_Run_CallsCorrectSubscription(t *testing.T) {
	c := make(MsgChan, 1)

	var callChan  = make(chan int)
	callback := func(b []byte) {
		callChan <- 1
	}

	c2 := make(MsgChan, 1)
	callback2 := func(b []byte) {
		callChan <- 2
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

	assertion: for {
		select {
		case c := <- callChan:
			if c == 2 {
				return
			}

			if c == 1 {
				t.Errorf("incorrect channel was called")
				return
			}
		case <- time.After(2 * time.Second):
			t.Errorf("failed asserting that any callbacks were called")
			break assertion
		}
	}
}
