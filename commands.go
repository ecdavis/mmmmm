package main

import (
	"errors"
	"strings"
)

type Command func(*User, string, []string)

var commandTable = make(map[string]Command)

func AddCommand(name string, command Command) {
	commandTable[name] = command
}

func RemoveCommand(name string) {
	delete(commandTable, name)
}

func ExistsCommand(name string) bool {
	_, ok := commandTable[name]
	return ok
}

func RunCommand(user *User, cmd string, args []string) error {
	command, ok := commandTable[cmd]
	if !ok {
		return errors.New("command does not exist")
	}
	command(user, cmd, args)
	return nil
}

func HandleCommand(game *Game, user *User, input string) {
	words := strings.Split(input, " ")
	if len(words) < 1 {
		// TODO error
	}
	cmd, args := words[0], words[1:]
	RunCommand(user, cmd, args)
}
