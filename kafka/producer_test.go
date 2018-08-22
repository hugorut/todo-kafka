package kafka

import (
	"testing"
	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"log"
	"github.com/hugorut/todo-kafka/test"
	"github.com/stretchr/testify/assert"
	"github.com/pkg/errors"
)

func TestNewProducer(t *testing.T) {
	seedBroker := sarama.NewMockBroker(t, 1)

	metadataResponse := new(sarama.MetadataResponse)
	seedBroker.Returns(metadataResponse)

	producer, err := NewProducer(seedBroker.Addr())

	if err != nil {
		t.Errorf("Failed asserting error was nil, instead was %+v", err)
		return
	}

	if _, ok := producer.(Producer); !ok {
		t.Error("Failed asserting that NewProducer returned sarama.Producer interface.")
	}
}

func TestSendMessage(t *testing.T) {
	w := &test.Writer{}
	log.SetOutput(w)

	p := mocks.NewSyncProducer(t, nil)
	p.ExpectSendMessageAndSucceed()

	SendMessage(p, &sarama.ProducerMessage{})
	assert.Contains(t, w.Msg(), "sent message to partition: 0 offset: 1")
}

func TestSendMessage_ReturnsErr(t *testing.T) {
	p := mocks.NewSyncProducer(t, nil)
	err := errors.New("my error")
	p.ExpectSendMessageAndFail(err)

	actual := SendMessage(p, &sarama.ProducerMessage{})
	assert.Equal(t, err, actual)
}
