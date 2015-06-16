package main

import (
	"fmt"
	"github.com/ecdavis/mmmmm/hooks"
	"github.com/ecdavis/mmmmm/net"
	"log"
	"strings"
)

func on_start(args ...interface{}) {
	println("on_start")
}

func think(user *User, cmd string, args []string) {
	user.session.Write <- fmt.Sprintf("You think, '%s'\r\n", strings.Join(args, " "))
}

func main() {
	AddCommand("think", think)
	hooks.Add("foo", on_start)

	game := NewGame()

	sessions, err := net.RunServer()
	if err != nil {
		log.Fatal("runServer:", err)
	}
	go func() {
		game.AddSession(<-sessions)
	}()

	hooks.Run("foo")
	game.ProcessTasks()
}
