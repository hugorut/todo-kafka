package main

import (
	"log"
	"github.com/hugorut/todo-kafka/kafka"
	"github.com/hugorut/todo-kafka/todo"
	"flag"
	"github.com/hugorut/todo-kafka/http"
)

var address = flag.String("address", "127.0.0.1:9092", "kafka service address")
var port = flag.String("http-port", "8081", "port that the http services should listen on")

func main() {
	flag.Parse()

	consumer, err := kafka.NewConsumer(*address)
	if err != nil {
		log.Fatal(err)
	}

	createMsgs, err := kafka.ConsumeTopic(consumer, kafka.Todos)
	if err != nil {
		log.Fatal(err)
	}

	b := NewBroker(
		NewSubscription(createMsgs, map[kafka.Event]todo.MessageCallback{
			kafka.CreateTodo: todo.Create,
			kafka.UpdateTodo: todo.Update,
			kafka.DeleteTodo: todo.Delete,
		}),
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
