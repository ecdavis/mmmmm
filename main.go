package main

import (
	"errors"
	"fmt"
	"github.com/ecdavis/mmmmm/commands"
	"github.com/ecdavis/mmmmm/game"
	"github.com/ecdavis/mmmmm/hooks"
	"github.com/ecdavis/mmmmm/net"
	"log"
	"strings"
)

func HandleCommand(game *game.Game, user *game.User, input string) {
	words := strings.Split(input, " ")
	if len(words) < 1 {
		// TODO error
	}
	cmd, args := words[0], words[1:]
	commands.Run(user, cmd, args)
}

func onAddUser(args ...interface{}) error {
	if len(args) != 1 {
		// TODO error
		return errors.New("wrong number of arguments to hook")
	}
	user := args[0].(*game.User)
	user.InputHandlers = append(user.InputHandlers, HandleCommand)
	return nil
}

func think(user *game.User, cmd string, args []string) error {
	user.Session.Write <- fmt.Sprintf("You think, '%s'\r\n", strings.Join(args, " "))
	return nil
}

func main() {
	commands.Add("think", think)
	hooks.Add("addUser", onAddUser)

	game := game.NewGame()

	sessions, err := net.RunServer()
	if err != nil {
		log.Fatal("runServer:", err)
	}
	go func() {
		game.AddSession(<-sessions)
	}()

	game.ProcessTasks()
}
