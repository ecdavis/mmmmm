package main

import (
	"log"
	"net"
)

var inputHandlerStack = make([]func(*SessionManager, *SessionInput), 0)

func handleInput(manager *SessionManager, sisChannel chan *SessionInput) {
	for si := range sisChannel {
		inputHandlerStack[len(inputHandlerStack)-1](manager, si)
	}
}

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
	ch := make(chan *SessionInput)

	inputHandlerStack = append(inputHandlerStack, HandleCommand)

	manager := NewSessionManager()
	go manager.ProcessCommands(ch)
	go handleInput(manager, ch)

	err := runServer(manager)
	if err != nil {
		log.Fatal("runServer:", err)
	}
}
