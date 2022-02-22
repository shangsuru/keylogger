package main

import (
	"bufio"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("Serving %s\n", remoteAddr)
	defer conn.Close()
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("%s disconnected\n", remoteAddr)
			break
		}
		log.Println(string(data))
	}
}

func main() {
	ln, err := net.Listen("tcp", ":2345")
	if err != nil {
		log.Fatalf("Could not create TCP server: %v\n", err)
	}
	defer ln.Close()
	log.Println("Listening on port 2345")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}
