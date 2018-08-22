package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/Shopify/sarama"
	"time"
)

func TestBroker_Run(t *testing.T) {
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

	b := NewBroker(NewSubscription(c, callback), NewSubscription(c2, callback2))
	b.Run()

	c2 <- &sarama.ConsumerMessage{}
	time.Sleep(1 * time.Second)
	assert.False(t, called)
	assert.True(t, called2)
}
