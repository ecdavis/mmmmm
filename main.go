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

func runServer() (<-chan *Session, error) {
	ln, err := net.Listen("tcp", ":4040")
	if err != nil {
		return nil, err
	}
	ch := make(chan *Session)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Print("Accept:", err)
				continue
			}
			ch <- NewSession(conn)
		}
		close(ch)
	}()
	return ch, nil
}

func main() {
	AddCommand("think", think)

	game := NewGame()

	sessions, err := runServer()
	if err != nil {
		log.Fatal("runServer:", err)
	}
	go func() {
		game.AddSession(<-sessions)
	}()

	game.ProcessTasks()
}
