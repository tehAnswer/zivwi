package main

import (
	"os"

	nsq "github.com/bitly/go-nsq"
)

type Queue interface {
	Publish(topic string, data []byte) error
}

func NewQueue() Queue {
	config := nsq.NewConfig()
	producer, _ := nsq.NewProducer(os.Getenv("NSQ_URL"), config)
	return producer
}
