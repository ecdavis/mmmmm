package main

import (
	"log"
	"net"
)

func runServer(manager *ClientManager) error {
	ln, err := net.Listen("tcp", ":4040")
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("Accept:", err)
		}
		client := NewClient(conn)
		manager.add <- client
	}
}

func main() {
	manager := NewClientManager()
	err := runServer(manager)
	if err != nil {
		log.Fatal("runServer:", err)
	}
}
