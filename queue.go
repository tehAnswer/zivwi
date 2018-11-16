package main

type Queue interface {
	Publish(topic string, data []byte) error
}

func NewQueue() Queue {
	return nil
}
