package main

import (
	"log"
	"net"
)

func runServer(manager *SessionManager) error {
	ln, err := net.Listen("tcp", ":4040")
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("Accept:", err)
		}
		session := NewSession(conn)
		manager.add <- session
	}
}

func main() {
	manager := NewSessionManager()

	ch := make(chan string)
	go func() {
		for {
			manager.write <- <-ch
		}
	}()
	go manager.ProcessCommands(ch)
	err := runServer(manager)
	if err != nil {
		log.Fatal("runServer:", err)
	}
}
