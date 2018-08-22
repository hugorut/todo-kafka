package kafka

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"time"
)

func TestNewConsumer(t *testing.T) {
	seedBroker := sarama.NewMockBroker(t, 1)

	metadataResponse := new(sarama.MetadataResponse)
	seedBroker.Returns(metadataResponse)

	c, err := NewConsumer(seedBroker.Addr())
	assert.Nil(t, err)

	if v, ok := c.(Consumer); !ok {
		t.Errorf("Failed asserting that NewConsumer returned Consumer interface, %v", v)
	}
}

func TestConsumeTopic(t *testing.T) {
	consumer := mocks.NewConsumer(t, nil)
	consumer.SetTopicMetadata(map[string][]int32{
		CreateToDo.String(): []int32{
			0,
		},
	})

	pc := consumer.ExpectConsumePartition(CreateToDo.String(), 0, sarama.OffsetNewest)

	msgs, err := ConsumeTopic(consumer, CreateToDo)
	assert.Nil(t, err)

	go func() {
		time.Sleep(1 * time.Second)
		pc.YieldMessage(&sarama.ConsumerMessage{
			Value: sarama.ByteEncoder([]byte(`hello`)),
		})
	}()

	assertion: for {
		select {
		case <- msgs:
			break assertion
		case <-time.After(3 * time.Second):
			t.Error("failed to receive message")
			break assertion
		}
	}

}
