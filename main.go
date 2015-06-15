package main

import (
	"log"
	"net"
)

var inputHandlerStack = make([]func(*Game, *SessionInput), 0)

func runServer(game *Game) error {
	ln, err := net.Listen("tcp", ":4040")
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("Accept:", err)
		}
		game.AddSession(NewSession(conn))
	}
}

func main() {
	inputHandlerStack = append(inputHandlerStack, HandleCommand)

	game := NewGame()
	go game.ProcessTasks()

	err := runServer(game)
	if err != nil {
		log.Fatal("runServer:", err)
	}
}
