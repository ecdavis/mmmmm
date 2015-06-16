package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func think(user *User, cmd string, args []string) {
	user.session.write <- fmt.Sprintf("You think, '%s'\r\n", strings.Join(args, " "))
}

var inputHandlerStack = make([]func(*Game, *User, string), 0)

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
	AddCommand("think", think)

	inputHandlerStack = append(inputHandlerStack, HandleCommand)

	game := NewGame()
	go game.ProcessTasks()

	err := runServer(game)
	if err != nil {
		log.Fatal("runServer:", err)
	}
}
