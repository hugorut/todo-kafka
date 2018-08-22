package main

import (
	"log"
	"github.com/hugorut/todo-kafka/kafka"
	"github.com/hugorut/todo-kafka/todo"
	"flag"
	"github.com/hugorut/todo-kafka/http"
)

var address = flag.String("address", "127.0.0.1:9092", "kafka service address")
var port = flag.String("http-port", "8081", "http listen port")

func main() {
	flag.Parse()

	consumer, err := kafka.NewConsumer(*address)
	if err != nil {
		log.Fatal(err)
	}

	createMsgs, err := kafka.ConsumeTopic(consumer, kafka.CreateToDo)
	if err != nil {
		log.Fatal(err)
	}

	deleteMsgs, err := kafka.ConsumeTopic(consumer, kafka.DeleteToDo)
	if err != nil {
		log.Fatal(err)
	}

	updateMsgs, err := kafka.ConsumeTopic(consumer, kafka.UpdateToDo)
	if err != nil {
		log.Fatal(err)
	}

	b := NewBroker(
		NewSubscription(createMsgs, todo.Create),
		NewSubscription(deleteMsgs, todo.Delete),
		NewSubscription(updateMsgs, todo.Update),
	)
	go b.Run()

	producer, err := kafka.NewProducer(*address)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(
		http.NewServer(producer, *port),
	)
}
