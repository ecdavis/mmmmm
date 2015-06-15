package main

import (
	"log"
	"net"
)

var inputHandlerStack = make([]func(*Game, *SessionInput), 0)

func handleInput(game *Game, sisChannel chan *SessionInput) {
	for si := range sisChannel {
		inputHandlerStack[len(inputHandlerStack)-1](game, si)
	}
}

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
		session := NewSession(conn)
		game.add <- session
	}
}

func main() {
	ch := make(chan *SessionInput)

	inputHandlerStack = append(inputHandlerStack, HandleCommand)

	game := NewGame()
	go game.ProcessCommands(ch)
	go handleInput(game, ch)

	err := runServer(game)
	if err != nil {
		log.Fatal("runServer:", err)
	}
}
