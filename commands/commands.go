package commands

import (
	"errors"
	"github.com/ecdavis/mmmmm/game"
)

type Command func(*game.User, string, []string) error

var commands = make(map[string]Command)

func Add(name string, command Command) error {
	_, ok := commands[name]
	if ok {
		return errors.New("command already exists")
	}
	commands[name] = command
	return nil
}

func Run(user *game.User, cmd string, args []string) error {
	command, ok := commands[cmd]
	if !ok {
		return errors.New("command does not exist")
	}
	return command(user, cmd, args)
}
