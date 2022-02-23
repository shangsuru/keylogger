package main

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func main() {
	amqpConn, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatal(err)
	}
	defer amqpConn.Close()
	channelAmqp, _ := amqpConn.Channel()
	forever := make(chan bool)
	msgs, _ := channelAmqp.Consume(
		os.Getenv("RABBITMQ_QUEUE"),
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("Waiting for messages")
	<-forever
}
