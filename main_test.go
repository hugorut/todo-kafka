package main

import (
	"testing"
	"github.com/Shopify/sarama"
	"flag"
	"github.com/hugorut/todo-kafka/kafka"
	"net/http"
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"time"
	"log"
	"os"
	"github.com/hugorut/todo-kafka/storage"
)

var testMsg = sarama.StringEncoder("Foo")

func Test_main_consumesTopics(t *testing.T) {
	t.Skip()

	sarama.Logger = log.New(os.Stdout, "sarama: ", 0)
	seedBroker := sarama.NewMockBroker(t, 1)

	metadataResponse := new(sarama.MetadataResponse)
	seedBroker.Returns(metadataResponse)

	seedBroker.SetHandlerByMap(map[string]sarama.MockResponse{
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(seedBroker.Addr(), seedBroker.BrokerID()).
			SetLeader(kafka.CreateToDo.String(), 0, seedBroker.BrokerID()).
			SetLeader(kafka.DeleteToDo.String(), 0, seedBroker.BrokerID()).
			SetLeader(kafka.UpdateToDo.String(), 0, seedBroker.BrokerID()),
		"OffsetRequest": sarama.NewMockOffsetResponse(t).
			SetOffset(kafka.CreateToDo.String(), 0, sarama.OffsetNewest, 10).
			SetOffset(kafka.CreateToDo.String(), 0, sarama.OffsetOldest, 7).
			SetOffset(kafka.DeleteToDo.String(), 0, sarama.OffsetNewest, 10).
			SetOffset(kafka.DeleteToDo.String(), 0, sarama.OffsetOldest, 7).
			SetOffset(kafka.UpdateToDo.String(), 0, sarama.OffsetNewest, 10).
			SetOffset(kafka.UpdateToDo.String(), 0, sarama.OffsetOldest, 7),
		"FetchRequest": sarama.NewMockFetchResponse(t, 1).
			SetMessage(kafka.CreateToDo.String(), 0, 9, testMsg).
			SetHighWaterMark(kafka.CreateToDo.String(), 0, 14).
			SetMessage(kafka.DeleteToDo.String(), 0, 9, testMsg).
			SetHighWaterMark(kafka.DeleteToDo.String(), 0, 14).
			SetMessage(kafka.UpdateToDo.String(), 0, 9, testMsg).
			SetHighWaterMark(kafka.UpdateToDo.String(), 0, 14),
	})

	flag.Set("address", seedBroker.Addr())
	flag.Set("http-port", "8082")

	go main()

	time.Sleep(2 * time.Second)
	buf := bytes.NewBufferString(`{"name":"test"}`)
	res, err := http.Post("http://127.0.0.1:8082/todo/create", "application/json", buf)
	if !assert.Nil(t, err) {
		return
	}

	b, _ := ioutil.ReadAll(res.Body)
	assert.JSONEq(t, `{"msg":"success"}`, string(b))

	assert.Len(t, storage.S.List(), 1)
	storage.S.Flush()
}
