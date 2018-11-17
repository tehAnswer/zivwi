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
	url := os.Getenv("NSQ_URL")
	producer, _ := nsq.NewProducer(url, config)
	return producer
}
