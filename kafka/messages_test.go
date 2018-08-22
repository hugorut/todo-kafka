package kafka

import (
	"testing"
	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func TestNewMessage(t *testing.T) {
	val := map[string]string{
		"msg": "hello",
	}

	msg := NewMessage("topic", 1, val)

	assert.IsType(t, &sarama.ProducerMessage{}, msg)
	assert.Equal(t, "topic", msg.Topic)
	assert.Equal(t, int32(1), msg.Partition)

	b, _ := msg.Value.Encode()
	var actual map[string]string
	json.Unmarshal(b, &actual)
	assert.EqualValues(t, val, actual)
}
