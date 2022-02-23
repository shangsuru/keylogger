package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func handleConnection(tcpConn net.Conn, amqpConn *amqp.Connection) {
	remoteAddr := tcpConn.RemoteAddr().String()
	log.Printf("Serving %s\n", remoteAddr)
	defer tcpConn.Close()

	channelAmqp, err := amqpConn.Channel()
	if err != nil {
		log.Printf("Could not connect to RabbitMQ: %v\n", err)
		return
	}

	for {
		data, err := bufio.NewReader(tcpConn).ReadString('\n')
		if err != nil {
			log.Printf("%s disconnected\n", remoteAddr)
			break
		}

		// Send data to RabbitMQ
		log.Printf("%s: \"%s\"\n", remoteAddr, strings.TrimSuffix(data, "\n"))
		payload := strings.Split(remoteAddr, ":")[0] + ":" + data // add the IP of the sender without the port to the data
		err = channelAmqp.Publish(
			"",
			os.Getenv("RABBITMQ_QUEUE"),
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				Timestamp:    time.Now(),
				ContentType:  "text/plain",
				Body:         []byte(payload),
			},
		)
		if err != nil {
			log.Printf("Could not publish keystrokes: %v\n", err)
			continue
		}
	}
}

func main() {
	// Receive keystrokes via TCP
	ln, err := net.Listen("tcp", ":2345")
	if err != nil {
		log.Fatalf("Could not create TCP server: %v\n", err)
	}
	defer ln.Close()
	log.Println("Listening on port 2345")

	// Keystrokes are published via RabbitMQ
	amqpConn, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to RabbitMQ server")

	for {
		tcpConn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleConnection(tcpConn, amqpConn)
	}
}
